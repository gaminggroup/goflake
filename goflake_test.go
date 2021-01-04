package goflake

import (
	"log"
	"sync"
	"testing"
)

func TestFlakeNodeId(t *testing.T) {
	err := SetNodeId(-1)
	if err == nil {
		log.Fatal("error initialising FlakeGen, it should not allow Node ID to be less than 0")
	}
	err = SetNodeId(1024)
	if err == nil {
		log.Fatal("error initialising FlakeGen, it should not allow Node ID to be greater than 1023")
	}
}

func TestFlakeRace(t *testing.T) {
	go func() {
		for i := 0; i < 10000; i++ {
			go func() {
				for i := 0; i < 10000; i++ {
					_ = SetNodeId(int64(i % 1024))
				}
			}()
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			flake, err := NextId()
			if err == nil {
				if flake == nil {
					t.Fatal("something went wrong")
				}
			}
		}
	}()
}

func TestFlakeDuplicates(t *testing.T) {
	var slice []int64
	queue := make(chan int64)
	var wg sync.WaitGroup
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			flake, err := NextId()
			if err == nil {
				queue <- flake.Int64()
			}
			wg.Done()
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	for id := range queue {
		slice = append(slice, id)
	}

	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				t.Fatal("found duplicate IDs")
			}
		}
	}
}

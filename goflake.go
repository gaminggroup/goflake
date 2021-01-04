package goflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const defaultEpoch = 1577836800000 // Wednesday, 1 January 2020 00:00:00 GMT
const timeMask = 0b11111111111111111111111111111111111111111
const nodeMask = 0b1111111111
const counterMask = 0b111111111111

type flake int64

var mutex sync.Mutex
var epoch time.Time
var nodeId int64
var counter int64
var lastTime int64

func init() {
	mutex = sync.Mutex{}
	epoch = time.Unix(defaultEpoch / 1000, (defaultEpoch % 1000) * 1000000)
	nodeId = 0
	counter = 0
	lastTime = 0
}

func SetNodeId(id int64) error {
	mutex.Lock()
	defer mutex.Unlock()
	if id < 0 || id >= 1024 {
		return errors.New(fmt.Sprintf("node ID (%d) must be in the range of 0 to 1023", id))
	}
	nodeId = id
	return nil
}

func NextId() (*flake, error) {
	mutex.Lock()
	defer mutex.Unlock()
	t := time.Since(epoch).Milliseconds()
	if t <= lastTime { // allow it to be less (just in case it is not using the monotonic clock)
		counter++
	} else {
		lastTime = t
		counter = 0
	}
	if counter >= 4096 {
		return nil, errors.New("unable to create Flake, the counter has overflowed")
	}
	id := ((timeMask & lastTime) << 22) | ((nodeMask & nodeId) << 12) | (counterMask & counter)
	flake := flake(id)
	return &flake, nil
}

func (f flake) Int64() int64 {
	return int64(f)
}

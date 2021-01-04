Usage:

```go
package main

import (
    "fmt"
    "github.com/gaminggroup/goflake"
)

func main() {
    _ = goflake.SetNodeId(int64(123)) // optional
    flake, err := goflake.NextId()
    if err != nil {
        fmt.Println(flake.Int64())
    }
}
```

pcduino
=======

pcduino.go: `go_environment` for pcDuino

see [godoc](https://godoc.org/github.com/conclave/pcduino)

- [hardware/core](https://godoc.org/github.com/conclave/pcduino/hardware/core)
- [hardware/sunxi](https://godoc.org/github.com/conclave/pcduino/hardware/sunxi)

## Program Template

```go
package main

import (
  . "github.com/conclave/pcduino/hardware/core"
)

func init() {
  Init()
  setup()
}

func main() {
  for {
    loop()
  }
}

// impl your own setup() and loop()
func setup() {

}

func loop() {

}
```

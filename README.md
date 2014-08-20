pcduino
=======

pcduino.go: `go_environment` for pcDuino [![Build Status](https://travis-ci.org/conclave/pcduino.svg)](https://travis-ci.org/conclave/pcduino)

see [godoc](https://godoc.org/github.com/conclave/pcduino)

- [core](https://godoc.org/github.com/conclave/pcduino/core)
- [sunxi](https://godoc.org/github.com/conclave/pcduino/sunxi)
- [lib/spi](https://godoc.org/github.com/conclave/pcduino/lib/spi)
- [lib/wire](https://godoc.org/github.com/conclave/pcduino/lib/wire)

## Program Template

```go
package main

import (
  . "github.com/conclave/pcduino/core"
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

pcduino
=======

pcduino.go: `go_environment` for pcDuino [![Build Status](https://travis-ci.org/conclave/pcduino.svg)](https://travis-ci.org/conclave/pcduino)

see [godoc](https://godoc.org/github.com/conclave/pcduino)

- [core](https://godoc.org/github.com/conclave/pcduino/core)
- [lib/i2c](https://godoc.org/github.com/conclave/pcduino/lib/i2c)
- [lib/spi](https://godoc.org/github.com/conclave/pcduino/lib/spi)
- [module/pcd8544](https://godoc.org/github.com/conclave/pcduino/module/pcd8544)

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

Or use `bin/pcduino-gen`:

```shell
pcduino-gen mypackage filename.go
```
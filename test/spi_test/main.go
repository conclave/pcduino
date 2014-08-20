package main

import (
	"flag"
	"fmt"
	"strconv"

	. "github.com/conclave/pcduino/core"
	. "github.com/conclave/pcduino/lib/spi"
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

var spi *SPIDevice
var deviceId = 0

func setup() {
	flag.Parse()
	if flag.NArg() >= 1 {
		deviceId, _ = strconv.Atoi(flag.Arg(0))
	}
	if deviceId > 1 {
		deviceId = 1
	}
	spi = NewSPIDevice(byte(deviceId))
	spi.Begin()
	spi.SetDataMode(SPI_MODE3)
	spi.SetBitOrder(MSBFIRST)
	spi.SetClockDivider(SPI_CLOCK_DIV16)
}

func loop() {
	//MSB first
	fmt.Printf("spi flash id = 0x%x\n", ReadSPIFlashID(spi))
	Delay(2000)
}

func ReadSPIFlashID(spi *SPIDevice) int {
	var cmd_rdid byte = 0x9F
	id := make([]byte, 3)
	spi.Transfer(cmd_rdid, SPI_CONTINUE)
	id[0] = spi.Transfer(0x00, SPI_CONTINUE)
	id[1] = spi.Transfer(0x00, SPI_CONTINUE)
	id[2] = spi.Transfer(0x00, SPI_LAST)
	var ret int = int(id[0]) << 8
	ret |= int(id[1])
	ret <<= 8
	ret |= int(id[2])
	return ret
}

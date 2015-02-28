package nRF24L

import (
	"github.com/conclave/pcduino/core"
)

func spiRead(reg byte) byte {
	var ret byte = 0
	core.DigitalWrite(CSN, 0)
	spiRw(reg)
	ret = spiRw(0)
	core.DigitalWrite(CSN, 1)
	return ret
}

func spiRw(reg byte) byte {
	for i := 0; i < 8; i++ {
		if reg&0x80 != 0 {
			core.DigitalWrite(MOSI, 1)
		} else {
			core.DigitalWrite(MOSI, 0)
		}
		core.DigitalWrite(SCK, 1)
		reg <<= 1
		if core.DigitalRead(MISO) == 1 {
			reg |= 1
		}
		core.DigitalWrite(SCK, 0)
	}
	return reg
}

func spiRwReg(reg, value byte) byte {
	core.DigitalWrite(CSN, 0)
	status := spiRw(reg)
	spiRw(value)
	core.DigitalWrite(CSN, 1)
	return status
}

func spiReadBuf(reg byte, buf []byte) byte {
	core.DigitalWrite(CSN, 0)
	status := spiRw(reg)
	for i := 0; i < len(buf); i++ {
		buf[i] = spiRw(0)
	}
	core.DigitalWrite(CSN, 1)
	return status
}

func spiWriteBuf(reg byte, buf []byte) byte {
	core.DigitalWrite(CSN, 0)
	status := spiRw(reg)
	for i := 0; i < len(buf); i++ {
		spiRw(buf[i])
	}
	core.DigitalWrite(CSN, 1)
	return status
}

package spi

import (
	"fmt"
	"syscall"
	"unsafe"

	. "github.com/conclave/pcduino/hardware/core"
	. "github.com/conclave/pcduino/hardware/sunxi"
)

const (
	SPI_CONTINUE     = 0
	SPI_LAST         = 1
	SPI_MODE0        = 0x00
	SPI_MODE1        = 0x01
	SPI_MODE2        = 0x02
	SPI_MODE3        = 0x03
	SPI_CLOCK_DIV1   = 0x00
	SPI_CLOCK_DIV2   = 0x01
	SPI_CLOCK_DIV4   = 0x02
	SPI_CLOCK_DIV8   = 0x03
	SPI_CLOCK_DIV16  = 0x04
	SPI_CLOCK_DIV32  = 0x05
	SPI_CLOCK_DIV64  = 0x06
	SPI_CLOCK_DIV128 = 0x07
)

const spi_name = "/dev/spidev0.0"
const spi1_name = "/dev/spidev1.0"
const spi2_name = "/dev/spidev2.0"
const bits_per_word = 8

type SPIDevice struct {
	deviceId byte
	fd       int
	speed    int
}

func NewSPIDevice(bus byte) *SPIDevice {
	return &SPIDevice{
		deviceId: bus,
		fd:       0,
		speed:    0,
	}
}

func (this *SPIDevice) Close() {
	this.End()
}

func (this *SPIDevice) Begin() {
	switch this.deviceId {
	case 0:
		Hw_PinMode(SPI_CS, IO_SPI_FUNC)
		Hw_PinMode(SPI_MOSI, IO_SPI_FUNC)
		Hw_PinMode(SPI_MISO, IO_SPI_FUNC)
		Hw_PinMode(SPI_CLK, IO_SPI_FUNC)
	case 1:
		Hw_PinMode(SPIEX_CS, IO_SPIEX_FUNC)
		Hw_PinMode(SPIEX_CS, IO_SPIEX_FUNC)
		Hw_PinMode(SPIEX_CS, IO_SPIEX_FUNC)
		Hw_PinMode(SPIEX_CS, IO_SPIEX_FUNC)
	}
	var err error
	if this.fd == 0 {
		switch this.deviceId {
		case 0:
			this.fd, err = syscall.Open(spi_name, syscall.O_RDWR|syscall.O_CLOEXEC, 0666)
		case 1:
			this.fd, err = syscall.Open(spi2_name, syscall.O_RDWR|syscall.O_CLOEXEC, 0666)
			if err != nil {
				this.fd, err = syscall.Open(spi1_name, syscall.O_RDWR|syscall.O_CLOEXEC, 0666)
			}
		}
	}
	if err != nil || this.fd <= 0 {
		panic("can't open spi device")
	}
	this.speed = 500000
	var default_mode int = 0
	if err = Ioctl(this.fd, SPI_IOC_RD_MODE(), uintptr(unsafe.Pointer(&default_mode))); err != nil {
		panic("can't get spi mode: " + err.Error())
	}
	var max_speed int = 0
	if err = Ioctl(this.fd, SPI_IOC_RD_MAX_SPEED_HZ(), uintptr(unsafe.Pointer(&max_speed))); err != nil {
		panic("can't get max speed: " + err.Error())
	}
	this.speed = max_speed
	fmt.Printf("init spi mode: 0x%x\n", default_mode)
	fmt.Printf("bits per word: %d\n", bits_per_word)
	fmt.Printf("max speed: %d Hz", this.speed)
}

func (this *SPIDevice) End() {
	if this.fd != 0 {
		syscall.Close(this.fd)
		this.fd = 0
	}
}

func (this *SPIDevice) SetBitOrder(bitOrder byte) {
	if this.fd == 0 {
		return
	}
	var order int = 0
	var mode int = 0
	var err error
	switch bitOrder {
	case LSBFIRST:
		order = 1
	case MSBFIRST:
		order = 0
	}
	if err = Ioctl(this.fd, SPI_IOC_WR_LSB_FIRST(), uintptr(unsafe.Pointer(&order))); err != nil {
		panic("can't set bits order: " + err.Error())
	}
	if err = Ioctl(this.fd, SPI_IOC_RD_MODE(), uintptr(unsafe.Pointer(&mode))); err != nil {
		panic("can't get spi mode: " + err.Error())
	}
	fmt.Printf("set bit order - spi mode: 0x%x\n", mode)
}

func (this *SPIDevice) SetDataMode(mode int) {
	if this.fd == 0 {
		return
	}
	var smode int = 0
	var err error
	if err = Ioctl(this.fd, SPI_IOC_RD_MODE(), uintptr(unsafe.Pointer(&smode))); err != nil {
		panic("can't get spi mode: " + err.Error())
	}
	smode &= ^0x3
	smode |= (mode & 0x3)
	if err = Ioctl(this.fd, SPI_IOC_WR_MODE(), uintptr(unsafe.Pointer(&smode))); err != nil {
		panic("can't set spi mode: " + err.Error())
	}
	if err = Ioctl(this.fd, SPI_IOC_RD_MODE(), uintptr(unsafe.Pointer(&smode))); err != nil {
		panic("can't get spi mode: " + err.Error())
	}
	fmt.Printf("set data mode - spi mode: 0x%x\n", smode)
}

func (this *SPIDevice) SetClockDivider(rate int) {
	switch rate {
	// case SPI_CLOCK_DIV1:
	case SPI_CLOCK_DIV2:
		this.speed /= 2
	case SPI_CLOCK_DIV4:
		this.speed /= 4
	case SPI_CLOCK_DIV8:
		this.speed /= 8
	case SPI_CLOCK_DIV16:
		this.speed /= 16
	case SPI_CLOCK_DIV32:
		this.speed /= 32
	case SPI_CLOCK_DIV64:
		this.speed /= 64
	case SPI_CLOCK_DIV128:
		this.speed /= 128
	}
}

func (this *SPIDevice) Transfer(value byte, mode byte) byte {
	if this.fd == 0 {
		return 0
	}
	var delay_usecs uint16
	switch mode {
	case SPI_CONTINUE:
		delay_usecs = 0
	case SPI_LAST:
		delay_usecs = 0xAA55
	}
	var tx byte = value
	var rx byte = 0
	transfer := SPI_IOC_Transfer{}
	transfer.TX_buf = uint64(uintptr(unsafe.Pointer(&tx)))
	transfer.RX_buf = uint64(uintptr(unsafe.Pointer(&rx)))
	transfer.Length = 1
	transfer.Speed_hz = uint32(this.speed)
	transfer.Delay_usecs = delay_usecs
	transfer.Bits_per_word = bits_per_word
	var err error
	if err = Ioctl(this.fd, SPI_IOC_MESSAGE(1), uintptr(unsafe.Pointer(&transfer))); err != nil {
		panic("can't send spi message: " + err.Error())
	}
	return tx
}

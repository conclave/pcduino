// TwoWire - TWI/I2C library for Wiring
package wire

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	. "github.com/conclave/pcduino/core"
)

const I2CCLOCK_CHANGE = 0x0740

const (
	I2C_RETRIES = 0x0701 /* number of times a device address should be polled when not acknowledging */
	I2C_TIMEOUT = 0x0702 /* set timeout in units of 10 ms */
	/* NOTE: Slave address is 7 or 10 bits, but 10-bit addresses
	 * are NOT supported! (due to code brokenness)
	 */
	I2C_SLAVE       = 0x0703 /* Use this slave address */
	I2C_SLAVE_FORCE = 0x0706 /* Use this slave address, even if it is already in use by a driver! */
	I2C_TENBIT      = 0x0704 /* 0 for 7 bit addrs, != 0 for 10 bit */
	I2C_FUNCS       = 0x0705 /* Get the adapter functionality mask */
	I2C_RDWR        = 0x0707 /* Combined R/W transfer (one STOP only) */
	I2C_PEC         = 0x0708 /* != 0 to use PEC with SMBus */
	I2C_SMBUS       = 0x0720 /* SMBus transfer */
) // from <linux/i2c-dev.h>

type TwoWire struct {
	I2CHandle     int
	rxBuffer      []byte
	rxBufferIndex int
	txBuffer      []byte
	txBufferIndex int
	transmitting  byte
}

func NewTwoWire() *TwoWire {
	r := TwoWire{}
	r.rxBuffer = make([]byte, 32)
	r.txBuffer = make([]byte, 32)
	return &r
}

func (this *TwoWire) Close() {
	if this.I2CHandle != 0 {
		syscall.Close(this.I2CHandle)
		this.I2CHandle = 0
	}
}

// NOT IMPLEMENTED
func (this *TwoWire) onRequest() {
}

// NOT IMPLEMENTED
func (this *TwoWire) onReceive(arg int) {
}

// NOT IMPLEMENTED
// behind the scenes function that is called when data is received
func (this *TwoWire) onRequestService() {
}

// NOT IMPLEMENTED
// behind the scenes function that is called when data is received
func (this *TwoWire) onReceiveService(a []byte, b int) {
}

func (this *TwoWire) Begin() {
	this.rxBufferIndex = 0
	this.txBufferIndex = 0
	if this.I2CHandle == 0 {
		var err error
		if this.I2CHandle, err = syscall.Open("/dev/i2c-2", syscall.O_RDWR|syscall.O_CLOEXEC, 0666); err != nil {
			panic("can't open i2c device: " + err.Error())
		}
	}
}

// bus freq range 10kHz-400kHz
func (this *TwoWire) SetBusFreq(speed_hz uint) {
	if speed_hz > 400000 || speed_hz < 10000 {
		fmt.Fprintf(os.Stderr, "invalid bus freq, range[10000,400000]; freq=%d\n", speed_hz)
		return
	}
	fd, err := syscall.Open("/dev/hwi2c", syscall.O_RDWR|syscall.O_CLOEXEC, 0666)
	if err != nil {
		panic("can't open i2c device: " + err.Error())
	}
	if err = Ioctl(fd, I2CCLOCK_CHANGE, uintptr(unsafe.Pointer(&speed_hz))); err != nil {
		panic("change i2c bus freq fail: " + err.Error())
	}
	syscall.Close(fd)
}

func (this *TwoWire) BeginTransmission(address int) {
	if this.I2CHandle == 0 {
		return
	}
	this.transmitting = 1
	this.txBufferIndex = 0
	var err error
	if address <= 0x7F {
		if err = Ioctl(this.I2CHandle, I2C_TENBIT, 0); err != nil {
			panic("set i2c 7-bits address flag fail: " + err.Error())
		}
	} else if address <= 0x3FF {
		if err = Ioctl(this.I2CHandle, I2C_TENBIT, 1); err != nil {
			panic("set i2c 10-bits address flag fail: " + err.Error())
		}
	}
	if err = Ioctl(this.I2CHandle, I2C_SLAVE, uintptr(address)); err != nil {
		panic("set i2c slave address fail: " + err.Error())
	}
}

func (this *TwoWire) EndTransmission() int {
	if this.I2CHandle == 0 {
		return 0
	}
	var ret int
	file := os.NewFile(uintptr(this.I2CHandle), "/dev/i2c-2")
	n, err := file.Write(this.txBuffer)
	if err != nil || n != len(this.txBuffer) {
		fmt.Fprintf(os.Stderr, "i2c transaction failed\n")
		ret = 4
	} else {
		ret = 0
	}
	this.txBufferIndex = 0
	this.transmitting = 0
	return ret
}

func (this *TwoWire) RequestFrom(address, quantity int) int {
	if this.I2CHandle == 0 {
		return 0
	}
	if quantity > 32 {
		quantity = 32
	}
	var err error
	if address <= 0x7F {
		if err = Ioctl(this.I2CHandle, I2C_TENBIT, 0); err != nil {
			panic("set i2c 7-bits address flag fail: " + err.Error())
		}
	} else if address <= 0x3FF {
		if err = Ioctl(this.I2CHandle, I2C_TENBIT, 1); err != nil {
			panic("set i2c 10-bits address flag fail: " + err.Error())
		}
	}
	if err = Ioctl(this.I2CHandle, I2C_SLAVE, uintptr(address)); err != nil {
		panic("set i2c slave address fail: " + err.Error())
	}
	this.rxBufferIndex = 0
	file := os.NewFile(uintptr(this.I2CHandle), "/dev/i2c-2")
	n, err := file.Read(this.rxBuffer[:quantity])
	if err != nil {
		return 0
	}
	this.rxBuffer = this.rxBuffer[:n]
	return n
}

// must be called in:
// slave tx event callback
// or after BeginTransmission(address)
func (this *TwoWire) Write(data []byte) int {
	if this.txBufferIndex >= 32 {
		return 0
	}
	ret := len(data)
	if this.transmitting == 1 {
		if this.txBufferIndex+ret > 32 {
			ret = 32 - this.txBufferIndex
			data = data[:ret]
		}
		copy(this.txBuffer[this.txBufferIndex:], data)
		this.txBufferIndex += ret
		this.txBuffer = this.txBuffer[:this.txBufferIndex]
	}
	return ret
}

// must be called in:
// slave rx event callback
// or after RequestFrom(address, numBytes)
func (this *TwoWire) Available() int {
	return len(this.rxBuffer) - this.rxBufferIndex
}

// must be called in:
// slave rx event callback
// or after RequestFrom(address, numBytes)
func (this *TwoWire) Read() int {
	value := -1
	if this.rxBufferIndex < len(this.rxBuffer) {
		value = int(this.rxBuffer[this.rxBufferIndex])
		this.rxBufferIndex++
	}
	return value
}

// must be called in:
// slave rx event callback
// or after RequestFrom(address, numBytes)
func (this *TwoWire) Peek() int {
	return 0
}

// NOT IMPLEMENTED
func (this *TwoWire) Flush() {
}

// NOT IMPLEMENTED
// sets function called on slave write
func (this *TwoWire) OnRequest(fn func(arg int)) {
}

// NOT IMPLEMENTED
// sets function called on slave read
func (this *TwoWire) OnReceive(fn func()) {
}

// Preinstantiate Objects
var Wire *TwoWire

func init() {
	Wire = NewTwoWire()
}

// Package i2c provides low level control over the linux i2c bus.
//
// Before usage you should load the i2c-dev kenel module
//
//      sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// linux systems contain several buses.
//
// based on https://github.com/davecheney/i2c/blob/master/i2c.go
package i2c

import (
	"fmt"
	"os"
	"syscall"
)

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

// I2C represents a connection to an i2c device.
type I2C struct {
	rc *os.File
}

// New opens a connection to an i2c device.
func New(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	return &I2C{f}, nil
}

// Write sends buf to the remote i2c device. The interpretation of
// the message is implementation dependant.
func (i2c *I2C) Write(buf ...byte) error {
	_, err := i2c.rc.Write(buf)
	return err
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

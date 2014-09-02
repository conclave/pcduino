// +build linux

package core

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const serial_name = "/dev/ttyS1"
const SERIAL_BUFFER_SIZE = 1024

type hwSerial struct {
	file  *os.File
	mutex *sync.Mutex
}

var Serial *hwSerial = nil

func (*hwSerial) Begin(baud uint, config byte) (err error) {
	if Serial == nil {
		f, err := os.OpenFile(serial_name, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
		if err != nil {
			return err
		}
		Serial = &hwSerial{f, &sync.Mutex{}}
		runtime.SetFinalizer(f, func(fd *os.File) {
			fd.Close()
		})
	}
	defer Serial.mutex.Unlock()
	Serial.mutex.Lock()
	Hw_PinMode(GPIO0, IO_UART_FUNC)
	Hw_PinMode(GPIO1, IO_UART_FUNC)

	fd := int(Serial.file.Fd())
	t := syscall.Termios{}
	if err = Ioctl(fd, syscall.TCGETS, uintptr(unsafe.Pointer(&t))); err != nil {
		return err
	}
	if err = Ioctl(fd, syscall.TCIOFLUSH, 0); err != nil {
		return err
	}
	rate := uint32(get_valid_baud(baud))
	if rate > 0 {
		t.Ispeed = rate
		t.Ospeed = rate
	}
	t.Cflag = uint32(int32(t.Cflag) & ^syscall.CSIZE)
	t.Cflag |= uint32(get_databit(config))

	switch get_parity(config) {
	case 'O':
		t.Cflag |= (syscall.PARODD | syscall.PARENB)
		t.Iflag |= syscall.INPCK
	case 'E':
		t.Cflag |= syscall.PARENB
		t.Cflag = uint32(int32(t.Cflag) & ^syscall.PARODD)
		t.Iflag |= syscall.INPCK
	default:
		t.Cflag = uint32(int32(t.Cflag) & ^syscall.PARENB)
		t.Iflag = uint32(int32(t.Iflag) & ^syscall.INPCK)
	}

	switch get_stopbit(config) {
	case 2:
		t.Cflag |= syscall.CSTOPB
	default:
		t.Cflag = uint32(int32(t.Cflag) & ^syscall.CSTOPB)
	}

	t.Cflag = uint32(int32(t.Cflag) & ^(syscall.ICANON | syscall.ECHO | syscall.ECHOE | syscall.ISIG))
	if err = Ioctl(fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&t))); err != nil {
		return err
	}
	if err = Ioctl(fd, syscall.TCIOFLUSH, 0); err != nil {
		return err
	}
	return nil
}

func (*hwSerial) Read(buffer []byte) (int, error) {
	defer Serial.mutex.Unlock()
	Serial.mutex.Lock()
	return Serial.file.Read(buffer)
}

func (*hwSerial) Flush() (err error) {
	defer Serial.mutex.Unlock()
	Serial.mutex.Lock()
	// if err = Serial.file.Sync(); err != nil {
	// 	return
	// }
	err = Ioctl(int(Serial.file.Fd()), syscall.TCIOFLUSH, 0)
	return
}

func (*hwSerial) Write(buffer []byte) (int, error) {
	defer Serial.mutex.Unlock()
	Serial.mutex.Lock()
	return Serial.file.Write(buffer)
}

func (*hwSerial) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(Serial, a...)
}

func (*hwSerial) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(Serial, format, a...)
}

func (*hwSerial) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(Serial, a...)
}

const (
	SERIAL_5N1 = 0x00
	SERIAL_6N1 = 0x02
	SERIAL_7N1 = 0x04
	SERIAL_8N1 = 0x06 //default
	SERIAL_5N2 = 0x08
	SERIAL_6N2 = 0x0A
	SERIAL_7N2 = 0x0C
	SERIAL_8N2 = 0x0E
	SERIAL_5E1 = 0x20
	SERIAL_6E1 = 0x22
	SERIAL_7E1 = 0x24
	SERIAL_8E1 = 0x26
	SERIAL_5E2 = 0x28
	SERIAL_6E2 = 0x2A
	SERIAL_7E2 = 0x2C
	SERIAL_8E2 = 0x2E
	SERIAL_5O1 = 0x30
	SERIAL_6O1 = 0x32
	SERIAL_7O1 = 0x34
	SERIAL_8O1 = 0x36
	SERIAL_5O2 = 0x38
	SERIAL_6O2 = 0x3A
	SERIAL_7O2 = 0x3C
	SERIAL_8O2 = 0x3E
)

func get_databit(config byte) byte {
	switch config {
	case SERIAL_5N1, SERIAL_5N2, SERIAL_5E1, SERIAL_5E2, SERIAL_5O1, SERIAL_5O2:
		return syscall.CS5
	case SERIAL_6N1, SERIAL_6N2, SERIAL_6E1, SERIAL_6E2, SERIAL_6O1, SERIAL_6O2:
		return syscall.CS6
	case SERIAL_7N1, SERIAL_7N2, SERIAL_7E1, SERIAL_7E2, SERIAL_7O1, SERIAL_7O2:
		return syscall.CS7
		// case SERIAL_8N1, SERIAL_8N2, SERIAL_8E1, SERIAL_8E2, SERIAL_8O1, SERIAL_8O2:
		// default:
		// 	return CS8
	}
	return syscall.CS8
}

func get_stopbit(config byte) byte {
	switch config {
	case SERIAL_5N2, SERIAL_6N2, SERIAL_7N2, SERIAL_8N2, SERIAL_5E2, SERIAL_6E2, SERIAL_7E2, SERIAL_8E2, SERIAL_5O2, SERIAL_6O2, SERIAL_7O2, SERIAL_8O2:
		return 2
		// case SERIAL_5N1, SERIAL_6N1, SERIAL_7N1, SERIAL_8N1, SERIAL_5E1, SERIAL_6E1, SERIAL_7E1, SERIAL_8E1, SERIAL_5O1, SERIAL_6O1, SERIAL_7O1, SERIAL_8O1:
		// default:
		// 	return 1
	}
	return 1
}

func get_parity(config byte) byte {
	switch config {
	// case SERIAL_5N1, SERIAL_5N2, SERIAL_6N1, SERIAL_6N2, SERIAL_7N1, SERIAL_7N2, SERIAL_8N1, SERIAL_8N2:
	// default:
	// 	return 'N'
	case SERIAL_5O1, SERIAL_5O2, SERIAL_6O1, SERIAL_6O2, SERIAL_7O1, SERIAL_7O2, SERIAL_8O1, SERIAL_8O2:
		return 'O'
	case SERIAL_5E1, SERIAL_5E2, SERIAL_6E1, SERIAL_6E2, SERIAL_7E1, SERIAL_7E2, SERIAL_8E1, SERIAL_8E2:
		return 'E'
	}
	return 'N'
}

func get_valid_baud(speed uint) uint {
	switch speed {
	case 300:
		return syscall.B300
	case 600:
		return syscall.B600
	case 1200:
		return syscall.B1200
	case 2400:
		return syscall.B2400
	case 4800:
		return syscall.B4800
	case 9600:
		return syscall.B9600
	case 19200:
		return syscall.B19200
	case 38400:
		return syscall.B38400
	case 57600:
		return syscall.B57600
	case 115200:
		return syscall.B115200
	}
	return 0
}

package core

import (
	// "fmt"
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
	mutex sync.Mutex
}

var Serial *hwSerial = nil

func (this *hwSerial) Begin(baud uint, config byte) error {
	if Serial == nil {
		f, err := os.OpenFile(serial_name, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
		if err != nil {
			return err
		}
		Serial = &hwSerial{f, sync.Mutex{}}
		runtime.SetFinalizer(f, func(fd *os.File) {
			fd.Close()
		})
	}
	Hw_PinMode(GPIO0, IO_UART_FUNC)
	Hw_PinMode(GPIO1, IO_UART_FUNC)
	// set attribute
	rate := uint32(get_valid_baud(baud))
	if rate == 0 {
		return nil
	}
	t := syscall.Termios{
		Iflag:  syscall.IGNPAR,
		Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | rate,
		Cc:     [32]uint8{syscall.VMIN: 1},
		Ispeed: rate,
		Ospeed: rate,
	}
	fd := this.file.Fd()
	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); errno != 0 {
		return errno
	}
	if err := syscall.SetNonblock(int(fd), false); err != nil {
		return err
	}
	return nil
}

func (this *hwSerial) Read(buffer []byte) (int, error) {
	return this.file.Read(buffer)
}

func (this *hwSerial) Flush() error {
	return this.file.Sync()
}

func (this *hwSerial) Write(buffer []byte) (int, error) {
	return this.file.Write(buffer)
}

const (
	SERIAL_5N1 = 0x00
	SERIAL_6N1 = 0x02
	SERIAL_7N1 = 0x04
	SERIAL_8N1 = 0x06
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
		return _CS5
	case SERIAL_6N1, SERIAL_6N2, SERIAL_6E1, SERIAL_6E2, SERIAL_6O1, SERIAL_6O2:
		return _CS6
	case SERIAL_7N1, SERIAL_7N2, SERIAL_7E1, SERIAL_7E2, SERIAL_7O1, SERIAL_7O2:
		return _CS7
		// case SERIAL_8N1, SERIAL_8N2, SERIAL_8E1, SERIAL_8E2, SERIAL_8O1, SERIAL_8O2:
		// default:
		// 	return CS8
	}
	return _CS8
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
		return _B300
	case 600:
		return _B600
	case 1200:
		return _B1200
	case 2400:
		return _B2400
	case 4800:
		return _B4800
	case 9600:
		return _B9600
	case 14400:
		return 0
	case 19200:
		return _B19200
	case 28800:
		return 0
	case 38400:
		return _B38400
	case 57600:
		return _B57600
	case 115200:
		return _B115200
	}
	return 0
}

const ( // from /usr/include/x86_64-linux-gnu/bits/termios.h
	_CS5 = 0000000
	_CS6 = 0000020
	_CS7 = 0000040
	_CS8 = 0000060

	_B0       = 0000000 /* hang up */
	_B50      = 0000001
	_B75      = 0000002
	_B110     = 0000003
	_B134     = 0000004
	_B150     = 0000005
	_B200     = 0000006
	_B300     = 0000007
	_B600     = 0000010
	_B1200    = 0000011
	_B1800    = 0000012
	_B2400    = 0000013
	_B4800    = 0000014
	_B9600    = 0000015
	_B19200   = 0000016
	_B38400   = 0000017
	_B57600   = 0010001
	_B115200  = 0010002
	_B230400  = 0010003
	_B460800  = 0010004
	_B500000  = 0010005
	_B576000  = 0010006
	_B921600  = 0010007
	_B1000000 = 0010010
	_B1152000 = 0010011
	_B1500000 = 0010012
	_B2000000 = 0010013
	_B2500000 = 0010014
	_B3000000 = 0010015
	_B3500000 = 0010016
	_B4000000 = 0010017
)

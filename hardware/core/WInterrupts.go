package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

const (
	SWIRQ_START   = 0x201
	SWIRQ_STOP    = 0x202
	SWIRQ_SETPID  = 0x203
	SWIRQ_ENABLE  = 0x204
	SWIRQ_DISABLE = 0x205
)

const swirq_dev = "/dev/swirq"

type SWIrq_Config struct {
	channel   uint8
	mode, pid int
}

type userFunc func()

var user1fn, user2fn userFunc
var c1, c2 chan os.Signal
var b1, b2 chan bool

func init() {
	c1 = make(chan os.Signal, 1)
	c2 = make(chan os.Signal, 1)
	b1 = make(chan bool, 1)
	b2 = make(chan bool, 1)
	signal.Notify(c1, syscall.SIGUSR1)
	signal.Notify(c2, syscall.SIGUSR2)
}

func AttachInterrupt(irqno uint8, fn userFunc, modec int) {
	hwmode := 0
	if irqno < EXTERNAL_NUM_INTERRUPTS && fn != nil {
		switch modec {
		case LOW:
			hwmode = 0x3
		case FALLING:
			hwmode = 0x1
		case RISING:
			hwmode = 0x0
		case CHANGE:
			hwmode = 0x4
		default:
			fmt.Fprintf(os.Stderr, "attachInterrupt error: set interrupt %d to invalid mode\n", irqno)
			return
		}
		switch irqno {
		case 0:
			if user1fn != nil {
				b1 <- true
			}
			user1fn = fn
			go func() {
				for {
					select {
					case <-c1:
						if user1fn != nil {
							user1fn()
						}
					case <-b1:
						user1fn = nil
						return
					}
				}
			}()
		case 1:
			if user2fn != nil {
				b2 <- true
			}
			user2fn = fn
			go func() {
				for {
					select {
					case <-c2:
						if user2fn != nil {
							user2fn()
						}
					case <-b2:
						user2fn = nil
						return
					}
				}
			}()
		}
		pid := os.Getpid()
		irqconfig := SWIrq_Config{irqno, hwmode, pid}
		fd, err := syscall.Open(swirq_dev, os.O_RDONLY, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open swirq device fail: %v\n", err)
			os.Exit(-1)
		}
		if err = Ioctl(fd, SWIRQ_STOP, uintptr(unsafe.Pointer(&irqno))); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			panic("can't set SWIRQ_STOP")
		}
		if err = Ioctl(fd, SWIRQ_SETPID, uintptr(unsafe.Pointer(&irqconfig))); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			panic("can't set SWIRQ_SETPID")
		}
		if err = Ioctl(fd, SWIRQ_START, uintptr(unsafe.Pointer(&irqno))); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			panic("can't set SWIRQ_START")
		}
		syscall.Close(fd)
	}
}

func DetachInterrupt(irqno uint8) {
	if irqno < EXTERNAL_NUM_INTERRUPTS {
		fd, err := syscall.Open(swirq_dev, os.O_RDONLY, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open swirq device fail: %v\n", err)
			os.Exit(-1)
		}
		if err = Ioctl(fd, SWIRQ_STOP, uintptr(unsafe.Pointer(&irqno))); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			panic("can't set SWIRQ_STOP")
		}
		syscall.Close(fd)
		switch irqno {
		case 0:
			b1 <- true
		case 1:
			b2 <- true
		}
	}
}

func Interrupts() {
	var irqno uint8 = 0
	fd, err := syscall.Open(swirq_dev, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open swirq device fail: %v\n", err)
		os.Exit(-1)
	}
	if err = Ioctl(fd, SWIRQ_ENABLE, uintptr(unsafe.Pointer(&irqno))); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		panic("can't set interrupt_0 SWIRQ_ENABLE")
	}
	irqno = 1
	if err = Ioctl(fd, SWIRQ_ENABLE, uintptr(unsafe.Pointer(&irqno))); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		panic("can't set interrupt_1 SWIRQ_ENABLE")
	}
	syscall.Close(fd)
}

func NoInterrupts() {
	var irqno uint8 = 0
	fd, err := syscall.Open(swirq_dev, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open swirq device fail: %v\n", err)
		os.Exit(-1)
	}
	if err = Ioctl(fd, SWIRQ_DISABLE, uintptr(unsafe.Pointer(&irqno))); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		panic("can't set interrupt_0 SWIRQ_DISABLE")
	}
	irqno = 1
	if err = Ioctl(fd, SWIRQ_DISABLE, uintptr(unsafe.Pointer(&irqno))); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		panic("can't set interrupt_1 SWIRQ_DISABLE")
	}
	syscall.Close(fd)
}

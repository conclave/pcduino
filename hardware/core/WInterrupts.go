package core

import (
	"fmt"
	"os"
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
} // TODO: Ioctl

type userFunc func()

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
		switch irqno { //TODO: impl signal
		case 0:
			// signal(SIGUSR1, (void (*) (int))userFunc);
		case 1:
			// signal(SIGUSR2, (void (*) (int))userFunc);
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

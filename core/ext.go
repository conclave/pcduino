package core

import (
	"fmt"
	"os"
	"runtime"
)

type Pin struct {
	pfd *os.File
	mfd *os.File
	io  byte
}

func (this *Pin) Init() (err error) {
	name := fmt.Sprintf("%s%s%d", GPIO_PIN_DIR, GPIO_IF_PREFIX, this.io)
	if this.pfd, err = os.OpenFile(name, os.O_RDWR, 0644); err != nil {
		return
	}
	name = fmt.Sprintf("%s%s%d", GPIO_MODE_DIR, GPIO_IF_PREFIX, this.io)
	if this.mfd, err = os.OpenFile(name, os.O_RDWR, 0644); err != nil {
		return
	}
	runtime.SetFinalizer(this, func(p *Pin) {
		if p.pfd != nil {
			p.pfd.Close()
		}
		if p.mfd != nil {
			p.mfd.Close()
		}
	})
	return nil
}

func NewPin(pin byte) (*Pin, error) {
	p := &Pin{}
	p.io = pin
	return p, p.Init()
}

func (this *Pin) Mode(mode byte) (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok {
			err = panicErr
		}
	}()
	PinMode(this.io, mode)
	return err
}

func (this *Pin) DigitalWrite(value byte) {
	DigitalWrite(this.io, value)
}

func (this *Pin) DigitalRead() byte {
	return DigitalRead(this.io)
}

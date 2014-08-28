package main

import (
	"fmt"

	. "github.com/conclave/pcduino/core"
)

const (
	MAX_COUNT = 30
	INT_MODE  = FALLING
)

var led0 byte = 18
var led1 byte = 19
var bc0, bc1 = 0, 0
var state0 byte = LOW
var state1 byte = LOW

func init() {
	Init()
	setup()
}

func main() {
	for {
		loop()
	}
}

func setup() {
	PinMode(led0, OUTPUT)
	PinMode(led1, OUTPUT)
	AttachInterrupt(0, blink0, INT_MODE)
	AttachInterrupt(1, blink1, INT_MODE)
}

func loop() {
	DigitalWrite(led0, state0)
	DigitalWrite(led1, state1)
	if bc0 >= MAX_COUNT {
		DetachInterrupt(0)
	}
	if bc1 >= MAX_COUNT {
		DetachInterrupt(1)
	}
	Delay(1000)
}

func blink0() {
	state0 = 1 - state0
	bc0++
	fmt.Printf("blink0: %d, count=%d\n", state0, bc0)
}

func blink1() {
	state1 = 1 - state1
	bc1++
	fmt.Printf("blink1: %d, count=%d\n", state1, bc1)
}

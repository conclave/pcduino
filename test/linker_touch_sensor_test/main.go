package main

import (
	. "github.com/conclave/pcduino/core"
)

func init() {
	Init()
	setup()
}

func main() {
	for {
		loop()
	}
}

var touchPin byte = 1
var ledPin byte = 0

func setup() {
	println("Touch sensor test code!")
	println("Using I/O_0=Drive LED, I/O_1=Sensor output.")
	PinMode(touchPin, INPUT)
	PinMode(ledPin, OUTPUT)
}

func loop() {
	val := DigitalRead(touchPin)
	DigitalWrite(ledPin, val)
}

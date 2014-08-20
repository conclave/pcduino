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

var pin byte = 0

func setup() {
	println("Relay test code!")
	println("Using I/O_0.")
	PinMode(pin, OUTPUT)
}

func loop() {
	DigitalWrite(pin, HIGH)
	Delay(5000)
	DigitalWrite(pin, LOW)
	Delay(5000)
}

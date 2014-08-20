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

var pin byte = 1

func setup() {
	PinMode(pin, OUTPUT)
}

func loop() {
	DigitalWrite(pin, HIGH) // set the LED on
	Delay(1000)             // wait for a second
	DigitalWrite(pin, LOW)  // set the LED off
	Delay(1000)             // wait for a second
}

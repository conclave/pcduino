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

var btnPin byte = 1
var ledPin byte = 0

func setup() {
	println("Butten Test Code!")
	println("Using I/O_0=Drive LED and I/O_1=Button input.")
	PinMode(ledPin, OUTPUT)
	PinMode(btnPin, INPUT)
}

func loop() {
	btnIn := DigitalRead(btnPin)
	DigitalWrite(ledPin, btnIn)
	Delay(10)
}

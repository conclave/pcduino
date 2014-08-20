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

var magneticPin byte = 1
var ledPin byte = 0

func setup() {
	println("Magnetic sensor test code!")
	println("Using I/O_0=Drive LED, I/O_1=Sensor output.")
	PinMode(magneticPin, INPUT)
	PinMode(ledPin, OUTPUT)
}

func loop() {
	value := DigitalRead(magneticPin)
	DigitalWrite(ledPin, value)
}

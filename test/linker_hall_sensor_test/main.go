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

var sensorPin byte = 2
var ledPin byte = 0

func setup() {
	println("Hall sensor test code!")
	println("Using I/O_0=Drive LED, I/O_2=Sensor output.")
	PinMode(sensorPin, INPUT)
	PinMode(ledPin, OUTPUT)
}

func loop() {
	value := DigitalRead(sensorPin)
	DigitalWrite(ledPin, value)
}

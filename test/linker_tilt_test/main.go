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

var ledPin byte = 0
var switchPin byte = 1
var val byte

func setup() {
	println("Tilt sensor test code!")
	println("Using I/O_0=Drive LED, I/O_1=Sensor output.")
	PinMode(ledPin, OUTPUT)
	PinMode(switchPin, INPUT)
}

func loop() {
	val = DigitalRead(switchPin)
	if val == HIGH {
		DigitalWrite(ledPin, HIGH)
	} else {
		DigitalWrite(ledPin, LOW)
	}
}

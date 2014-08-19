// Sample code for MQ2 Smoke Sensor Shield for pcDuino
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

var value int = 0
var count int = 0

func setup() {
	PinMode(7, OUTPUT)
}

func loop() {
	count++
	value = AnalogRead(0)
	if count == 3000 {
		count = 0
		println("sensor =", value)
	}
}

package main

import (
	"fmt"

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

func setup() {
	println("Joystick Test Code!")
	println("Using ADC_2 and ADC_3.")
}

func loop() {
	sensorValue := AnalogRead(2)
	sensorValue2 := AnalogRead(3)
	fmt.Printf("The X and Y coordinate is:%d x %d\n", sensorValue, sensorValue2)
	Delay(500)
}

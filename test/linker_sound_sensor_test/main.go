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

var adcPin byte = 0 // select the input pin for the potentiometer
var ledPin byte = 0 // select the pin for the LED
var adcIn int = 0   // variable to store the value coming from the sensor

func setup() {
	println("Sound sensor Test Code!")
	println("Using ADC_0 and I/O_0.")
	PinMode(ledPin, OUTPUT) // set ledPin to OUTPUT
}

func loop() {
	adcIn = AnalogRead(adcPin) // read the value from the sensor.
	if adcIn >= 50 {
		DigitalWrite(ledPin, HIGH) // if adc in >= 50, led light
	} else {
		DigitalWrite(ledPin, LOW)
	}
	fmt.Printf("adc:%d!\n", adcIn)
	Delay(20)
}

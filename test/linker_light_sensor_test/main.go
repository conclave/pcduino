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

var ledPin byte = 0
var thresholdValue = 10
var adcIn byte = 0

var filter []int
var sampleCount int

func setup() {
	println("Light sensor test code!")
	println("Using I/O_0=Drive LED, ADC_0=Sensor output.")
	PinMode(ledPin, OUTPUT)
	filter = make([]int, 4)
	filter[0] = AnalogRead(adcIn)
	filter[1] = AnalogRead(adcIn)
	filter[2] = AnalogRead(adcIn)
	filter[3] = AnalogRead(adcIn)
	sampleCount = 0
}

func loop() {
	sensorValue := AnalogRead(adcIn)
	var Rsensor int

	if sampleCount >= 3 {
		sampleCount = 0
	} else {
		sampleCount++
	}
	fmt.Printf("adc %d :%d!\n", sampleCount, sensorValue)
	filter[sampleCount] = sensorValue
	sensorValue = (filter[0] + filter[1] + filter[2] + filter[3]) / 4
	Rsensor = (sensorValue * 100) / 64
	if Rsensor < thresholdValue {
		DigitalWrite(ledPin, HIGH)
	} else {
		DigitalWrite(ledPin, LOW)
	}
	fmt.Printf("adc:%d!\n", sensorValue)
	Delay(200)
}

// Temperature sensor test program
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

//The analog pin the TMP36's Vout (sense) pin is connected to
//the resolution is 10 mV / degree centigrade with a
//500 mV offset to allow for negative temperatures

const sensorADC = 0

func setup() {
	println("Temperature sensor test code!")
	println("Using ADC_0=Sensor output.")
}

func loop() {
	//getting the voltage reading from the temperature sensor
	reading := AnalogRead(sensorADC)
	// converting that reading to voltage
	var voltage float64 = float64(reading) * 2.0
	voltage /= 64.0
	fmt.Printf("adc:%d\n", reading)
	fmt.Printf("%.2f volts\n", voltage)
	// now print out the temperature
	var temperatureC float64 = (voltage - 0.5) * 100 //converting from 10 mv per degree wit 500 mV offset
	//to degrees ((volatge - 500mV) times 100)
	fmt.Printf("%.2f degrees C\n", temperatureC)
	// now convert to Fahrenheight
	var temperatureF float64 = (temperatureC * 9.0 / 5.0) + 32.0
	fmt.Printf("%.2f degrees F\n\n", temperatureF)
	Delay(1000) //waiting a second
}

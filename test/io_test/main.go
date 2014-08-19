// I/O test program
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

var led_pin byte = 1
var btn_pin byte = 5

func setup() {
	fmt.Printf("press button (connected to pin %d) to turn on LED (connected to pin %d)\n", btn_pin, led_pin)
	PinMode(led_pin, OUTPUT)
	PinMode(btn_pin, INPUT)
}

func loop() {
	value := DigitalRead(btn_pin) // get button status
	if value == HIGH {            // button pressed
		DigitalWrite(led_pin, HIGH) // turn on LED
	} else { // button released
		DigitalWrite(led_pin, LOW) // turn off LED
	}
	Delay(100)
}

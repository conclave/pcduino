// PWM test program
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

var pwm_id byte = 5
var freq uint = 781
var value int = MAX_PWM_LEVEL / 2

func setup() {
	step := PWMFreqSet(pwm_id, freq)
	fmt.Printf("PWM%d set freq %d and valid duty cycle range [0, %d]\n", pwm_id, freq, step)
	if step > 0 {
		fmt.Printf("PWM%d test with duty cycle %d\n", pwm_id, value)
		AnalogWrite(pwm_id, value)
	}
}

func loop() {
	DelayMicroseconds(200000)
}

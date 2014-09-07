// ADC test program
package main

import (
	"flag"
	"fmt"
	"strconv"

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

var adc_id int = 0

func setup() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: program ADC_ID(0/1/2/3/4/5/6/7/8/9/10/11)\n")
		fmt.Println("Default will get ADC0 value")
	}
	adc_id, _ = strconv.Atoi(flag.Arg(0))
}

func loop() {
	value := AnalogRead(byte(adc_id))
	fmt.Printf("ADC%d level is %d\n", adc_id, value)
	DelayMicroseconds(100000)
}

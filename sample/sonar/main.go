package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	. "github.com/conclave/pcduino/core"
)

const (
	SPEED_OF_SOUND float64 = 343.0 / 2 / 10000
)

var trig byte
var echo byte

func init() {
	Init()
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "invoke: %v\n", os.ErrInvalid)
		os.Exit(-1)
	}
	trig_, _ := strconv.Atoi(flag.Arg(0))
	echo_, _ := strconv.Atoi(flag.Arg(1))
	trig = byte(trig_)
	echo = byte(echo_)
	PinMode(echo, INPUT)
	PinMode(trig, OUTPUT)
	DigitalWrite(trig, LOW)
	Delay(20)
}

func loop() {
	DigitalWrite(trig, HIGH)
	DelayMicroseconds(20)
	DigitalWrite(trig, LOW)
	duration := PulseIn(echo, HIGH, 1000000)
	fmt.Printf("cm: %.2f\n", float64(duration)*SPEED_OF_SOUND)
	Delay(80)
}

func main() {
	for {
		loop()
	}
}

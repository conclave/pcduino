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

var dataPin byte = 0
var clkPin byte = 1
var clkFlag int = 0

const mode = 0
const ON = 0xff
const OFF = 0x00

func setup() {
	println("LED bar test code!")
	println("Using I/O_0=DATA, I/O_1=CLK.")
	PinMode(dataPin, OUTPUT)
	PinMode(clkPin, OUTPUT)
	DigitalWrite(dataPin, LOW)
	DigitalWrite(clkPin, LOW)
	clkFlag = 0
}

func loop() {
	send16bitData(mode)
	sendLED(0x0155)
	latchData()
	Delay(2000)
	send16bitData(mode)
	sendLED(0x02AA)
	latchData()
	Delay(2000)
}

func send16bitData(data uint) {
	for i := 0; i < 16; i++ {
		if data&0x8000 != 0 {
			DigitalWrite(dataPin, HIGH)
		} else {
			DigitalWrite(dataPin, LOW)
		}
		if clkFlag == 1 {
			DigitalWrite(clkPin, LOW)
			clkFlag = 0
		} else {
			DigitalWrite(clkPin, HIGH)
			clkFlag = 1
		}
		data <<= 1
	}
}

func latchData() {
	latchFlag := 0
	DigitalWrite(dataPin, LOW)
	DelayMicroseconds(200)
	for i := 0; i < 8; i++ {
		if latchFlag == 1 {
			DigitalWrite(dataPin, LOW)
			latchFlag = 0
		} else {
			DigitalWrite(dataPin, HIGH)
			latchFlag = 1
		}
	}
	DelayMicroseconds(200)
}

func sendLED(state uint) {
	for i := 0; i < 12; i++ {
		if state&0x01 != 0 {
			send16bitData(ON)
		} else {
			send16bitData(OFF)
		}
		state >>= 1
	}
}

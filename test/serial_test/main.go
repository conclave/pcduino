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

var buffer []byte

func setup() {
	Serial.Begin(115200, SERIAL_8N1)
	for Serial == nil {
		Delay(10)
	}
	Serial.Println("Serial online")
	buffer = make([]byte, 1)
}

func loop() {
	if n, _ := Serial.Read(buffer); n > 0 {
		Serial.Printf("Received: %d", buffer[0])
	}
	Delay(200)
}

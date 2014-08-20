package main

import (
	"time"

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

var speakerPin byte = 1
var length int = 15
var notes = []byte{'c', 'c', 'g', 'g', 'a', 'a', 'g', 'f', 'f', 'e', 'e', 'd', 'd', 'c', ' '}
var beats = []byte{1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 2, 4}
var tempo int = 300

func setup() {
	println("Buzzer test code!")
	println("Using I/O_1=D1, I/O_2=D2.")
	PinMode(speakerPin, OUTPUT)
	PinMode(2, OUTPUT)
	DigitalWrite(2, LOW)
}

func loop() {
	for i := 0; i < length; i++ {
		if notes[i] == ' ' {
			Delay(time.Duration(int(beats[i]) * tempo))
		} else {
			PlayNote(notes[i], int64(int(beats[i])*tempo))
		}
		Delay(time.Duration(tempo / 2))
	}
}

func playTone(tone int64, duration int64) {
	var i int64
	for i = 0; i < duration*1000; i += tone * 2 {
		DigitalWrite(speakerPin, HIGH)
		DelayMicrosends(time.Duration(tone))
		DigitalWrite(speakerPin, LOW)
		DelayMicrosends(time.Duration(tone))
	}
}

func PlayNote(note byte, duration int64) {
	names := []byte{'c', 'd', 'e', 'f', 'g', 'a', 'b', 'C'}
	tones := []int64{1915, 1700, 1519, 1432, 1275, 1136, 1014, 956}
	for i := 0; i < 8; i++ {
		if names[i] == note {
			playTone(tones[i], duration)
		}
	}
}

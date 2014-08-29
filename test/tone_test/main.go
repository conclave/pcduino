package main

import (
	"time"

	. "github.com/conclave/pcduino/core"
)

const TONE_PIN = 5

var melody = []uint{NOTE_C4, NOTE_G3, NOTE_G3, NOTE_A3, NOTE_G3, 0, NOTE_B3, NOTE_C4}

// note durations: 4 = quarter note, 8 = eighth note, etc.:
var noteDurations = []uint{4, 8, 8, 4, 4, 4, 4, 4}

func init() {
	Init()
	setup()
}

func main() {
	for {
		loop()
	}
}

func setup() {
	noteDuration := uint(0)
	thisNote := 0
	for thisNote = 0; thisNote < 8; thisNote++ {
		// to calculate the note duration, take one second
		// divided by the note type.
		//e.g. quarter note = 1000 / 4, eighth note = 1000/8, etc.
		noteDuration = 1000 / noteDurations[thisNote]
		Tone(TONE_PIN, melody[thisNote])
		//pauseBetweenNotes = noteDuration * 1.30;
		Delay(time.Duration(noteDuration))
		// stop the tone playing:
		NoTone(TONE_PIN)
	}
}

func loop() {
	Tone(TONE_PIN, NOTE_A4)
	Delay(200)
	Tone(TONE_PIN, NOTE_B4)
	Delay(500)
	Tone(TONE_PIN, NOTE_C5)
	Delay(300)
}

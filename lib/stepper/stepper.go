package stepper

import (
	. "github.com/conclave/pcduino/core"
)

type Stepper struct {
	direction int
	speed     uint
	steps     uint
	delay     uint
	pins      []byte
	count     uint
	timestamp int64
}

func New(steps uint, pins ...byte) *Stepper {
	l := len(pins)
	if l != 2 && l != 4 {
		return nil
	}
	stepper := Stepper{
		direction: 0,
		speed:     0,
		steps:     steps,
		delay:     0,
		pins:      nil,
		count:     0,
		timestamp: 0,
	}
	if l == 2 {
		stepper.pins = []byte{pins[0], pins[1]}
		PinMode(pins[0], OUTPUT)
		PinMode(pins[1], OUTPUT)
	} else {
		stepper.pins = []byte{pins[0], pins[1], pins[2], pins[3]}
		PinMode(pins[0], OUTPUT)
		PinMode(pins[1], OUTPUT)
		PinMode(pins[2], OUTPUT)
		PinMode(pins[3], OUTPUT)
	}
	return &stepper
}

/*
  Drives a unipolar or bipolar stepper motor using  2 wires or 4 wires

  When wiring multiple stepper motors to a microcontroller,
  you quickly run out of output pins, with each motor requiring 4 connections.

  By making use of the fact that at any time two of the four motor
  coils are the inverse  of the other two, the number of
  control connections can be reduced from 4 to 2.

  A slightly modified circuit around a Darlington transistor array or an L293 H-bridge
  connects to only 2 microcontroler pins, inverts the signals received,
  and delivers the 4 (2 plus 2 inverted ones) output signals required
  for driving a stepper motor.

  The sequence of control signals for 4 control wires is as follows:

  Step C0 C1 C2 C3
     1  1  0  1  0
     2  0  1  1  0
     3  0  1  0  1
     4  1  0  0  1

  The sequence of controls signals for 2 control wires is as follows
  (columns C1 and C2 from above):

  Step C0 C1
     1  0  1
     2  1  1
     3  1  0
     4  0  0

  The circuits can be found at
  http://www.arduino.cc/en/Tutorial/Stepper
*/

func (this *Stepper) SetSpeed(speed uint) {
	this.delay = 60 * 1000 / this.steps / speed
}

func (this *Stepper) Step(n int) {
	nn := n
	if nn < 0 {
		nn = -nn
		this.direction = 0
	} else {
		this.direction = 1
	}
	for nn > 0 {
		if Millis()-this.timestamp >= int64(this.delay) {
			this.timestamp = Millis()
			if this.direction == 1 {
				this.count++
				if this.count == this.steps {
					this.count = 0
				}
			} else {
				if this.count == 0 {
					this.count = this.steps
				}
				this.count--
			}
			nn--
			this.step(byte(this.count % 4))
		}
	}
}

func (this *Stepper) Version() int {
	return 4
}

func (this *Stepper) step(n byte) {
	if len(this.pins) == 2 {
		switch n & 0x03 {
		case 0: // 01
			DigitalWrite(this.pins[0], LOW)
			DigitalWrite(this.pins[1], HIGH)
		case 1: // 11
			DigitalWrite(this.pins[0], HIGH)
			DigitalWrite(this.pins[1], HIGH)
		case 2: // 10
			DigitalWrite(this.pins[0], HIGH)
			DigitalWrite(this.pins[1], LOW)
		case 3: // 00
			DigitalWrite(this.pins[0], LOW)
			DigitalWrite(this.pins[1], LOW)
		}
	} else {
		switch n & 0x03 {
		case 0: // 1010
			DigitalWrite(this.pins[0], HIGH)
			DigitalWrite(this.pins[1], LOW)
			DigitalWrite(this.pins[2], HIGH)
			DigitalWrite(this.pins[3], LOW)
		case 1: // 0110
			DigitalWrite(this.pins[0], LOW)
			DigitalWrite(this.pins[1], HIGH)
			DigitalWrite(this.pins[2], HIGH)
			DigitalWrite(this.pins[3], LOW)
		case 2: // 0101
			DigitalWrite(this.pins[0], LOW)
			DigitalWrite(this.pins[1], HIGH)
			DigitalWrite(this.pins[2], LOW)
			DigitalWrite(this.pins[3], HIGH)
		case 3: // 1001
			DigitalWrite(this.pins[0], HIGH)
			DigitalWrite(this.pins[1], LOW)
			DigitalWrite(this.pins[2], LOW)
			DigitalWrite(this.pins[3], HIGH)
		}
	}
}

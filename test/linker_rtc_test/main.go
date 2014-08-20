package main

import (
	"flag"
	"fmt"
	"strconv"

	. "github.com/conclave/pcduino/core"
	. "github.com/conclave/pcduino/lib/wire"
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

const DS1307_I2C_ADDRESS = 0x68 // This is the I2C address
var wire *TwoWire

func setup() {
	flag.Parse()
	wire = NewTwoWire()
	wire.Begin()
	getDateDS1307()
	var second, minute, hour, dayOfWeek, dayOfMonth, month, year int
	if flag.NArg() > 1 {
		year, _ = strconv.Atoi(flag.Arg(0))
	}
	if flag.NArg() > 2 {
		month, _ = strconv.Atoi(flag.Arg(1))
	}
	if flag.NArg() > 3 {
		dayOfMonth, _ = strconv.Atoi(flag.Arg(2))
	}
	if flag.NArg() > 4 {
		hour, _ = strconv.Atoi(flag.Arg(3))
	}
	if flag.NArg() > 5 {
		minute, _ = strconv.Atoi(flag.Arg(4))
	}
	if flag.NArg() > 6 {
		second, _ = strconv.Atoi(flag.Arg(5))
	}
	//force setting
	setDateDS1307(byte(second), byte(minute), byte(hour), byte(dayOfWeek), byte(dayOfMonth), byte(month), byte(year))
}

func loop() {
	Delay(2000)
	getDateDS1307()
}

func decToBcd(val byte) byte {
	return ((val / 10 * 16) + (val % 10))
}

func bcdToDec(val byte) byte {
	return ((val / 16 * 10) + (val % 16))
}

func setDateDS1307(second, minute, hour, dayOfWeek, dayOfMonth, month, year byte) {
	wire.BeginTransmission(DS1307_I2C_ADDRESS)
	wire.Write([]byte{0})
	wire.Write([]byte{decToBcd(second)}) // 0 to bit 7 starts the clock
	wire.Write([]byte{decToBcd(minute)})
	wire.Write([]byte{decToBcd(hour)}) // If you want 12 hour am/pm you need to set
	// bit 6 (also need to change readDateDs1307)
	wire.Write([]byte{decToBcd(dayOfWeek)})
	wire.Write([]byte{decToBcd(dayOfMonth)})
	wire.Write([]byte{decToBcd(month)})
	wire.Write([]byte{decToBcd(year)})
	wire.EndTransmission()
}

func getDateDS1307() {
	wire.BeginTransmission(DS1307_I2C_ADDRESS)
	wire.Write([]byte{0})
	wire.EndTransmission()
	wire.RequestFrom(DS1307_I2C_ADDRESS, 7)
	// A few of these need masks because certain bits are control bits
	second := bcdToDec(byte(wire.Read() & 0x7f))
	minute := bcdToDec(byte(wire.Read()))
	hour := bcdToDec(byte(wire.Read() & 0x3f)) // Need to change this if 12 hour am/pm
	dayOfWeek := bcdToDec(byte(wire.Read()))
	dayOfMonth := bcdToDec(byte(wire.Read()))
	month := bcdToDec(byte(wire.Read()))
	year := bcdToDec(byte(wire.Read()))
	fmt.Printf("%d:%d:%d %d %d/%d/%d.\n", hour, minute, second, dayOfWeek, month, dayOfMonth, year)
}

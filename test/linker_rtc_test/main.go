package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	. "github.com/conclave/pcduino/core"
	. "github.com/conclave/pcduino/lib/i2c"
)

func init() {
	Init()
	setup()
}

func main() {
	if readOnce {
		return
	}
	for {
		loop()
	}
}

const DS1307_I2C_ADDRESS = 0x68 // This is the I2C address
var i2c *I2C
var readOnce bool

func setup() {
	var setTime bool
	flag.BoolVar(&setTime, "s", false, "set date")
	flag.BoolVar(&readOnce, "once", false, "read once")
	flag.Parse()
	var err error
	if i2c, err = New(DS1307_I2C_ADDRESS, 2); err != nil {
		panic(err.Error())
	}
	getDateDS1307()
	if setTime {
		var second, minute, hour, dayOfWeek, dayOfMonth, month, year int
		if flag.NArg() >= 1 {
			year, _ = strconv.Atoi(flag.Arg(0))
		}
		if flag.NArg() >= 2 {
			month, _ = strconv.Atoi(flag.Arg(1))
		}
		if flag.NArg() >= 3 {
			dayOfMonth, _ = strconv.Atoi(flag.Arg(2))
		}
		if flag.NArg() >= 4 {
			hour, _ = strconv.Atoi(flag.Arg(3))
		}
		if flag.NArg() >= 5 {
			minute, _ = strconv.Atoi(flag.Arg(4))
		}
		if flag.NArg() >= 6 {
			second, _ = strconv.Atoi(flag.Arg(5))
		}
		fmt.Printf("set date - %d:%d:%d %d/%d/%d\n", hour, minute, second, month, dayOfMonth, year)
		setDateDS1307(byte(second), byte(minute), byte(hour), byte(dayOfWeek), byte(dayOfMonth), byte(month), byte(year))
	}
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
	// 0 to bit 7 starts the clock
	// If you want 12 hour am/pm you need to set bit 6 of (hour) (also need to change getDateDS1307)
	i2c.Write(0,
		decToBcd(second),
		decToBcd(minute),
		decToBcd(hour),
		decToBcd(dayOfWeek),
		decToBcd(dayOfMonth),
		decToBcd(month),
		decToBcd(year))
}

func getDateDS1307() {
	i2c.Write(0)
	b := make([]byte, 7)
	if err := i2c.Read(b); err != nil {
		fmt.Fprintf(os.Stderr, "getDateDS1307: %v\n", err)
		return
	}
	// A few of these need masks because certain bits are control bits
	second := bcdToDec(b[0] & 0x7f)
	minute := bcdToDec(b[1])
	hour := bcdToDec(b[2] & 0x3f) // Need to change this if 12 hour am/pm
	_ = bcdToDec(b[3])
	dayOfMonth := bcdToDec(b[4])
	month := bcdToDec(b[5])
	year := bcdToDec(b[6])
	// fmt.Printf("%% %d:%d:%d %d/%d/%d\n", hour, minute, second, month, dayOfMonth, year)
	// pattern => 2013-11-19 15:11:40
	fmt.Printf("20%d-%d-%d %d:%d:%d\n", year, month, dayOfMonth, hour, minute, second)
}

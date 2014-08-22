// +build arm

package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	. "github.com/conclave/pcduino/lib/i2c"
)

func init() {
	setup()
}

func main() {
	if !readTime && !setTime && !writeTime {
		flag.Usage()
		return
	}
	if readTime {
		if t := getDateDS1307(); t != nil {
			fmt.Println(t)
		}
	}
	if setTime {
		if t := getDateDS1307(); t != nil {
			tv := syscall.Timeval{
				Sec:  int32(t.Unix()),
				Usec: int32(t.UnixNano() % 100000000),
			}
			if err := syscall.Settimeofday(&tv); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
	if writeTime {
		now := time.Now()
		fmt.Println("set time -", now)
		setDateDS1307(&now)
		time.Sleep(100 * time.Millisecond)
		if t := getDateDS1307(); t != nil {
			fmt.Println(t)
		}
	}
}

const DS1307_I2C_ADDRESS = 0x68 // This is the I2C address
var i2c *I2C
var readTime bool
var setTime bool
var writeTime bool

func setup() {
	flag.BoolVar(&readTime, "r", false, "read hardware clock and print result")
	flag.BoolVar(&setTime, "s", false, "set the system time from the hardware clock")
	flag.BoolVar(&writeTime, "w", false, "set the hardware clock from the current system time")
	flag.Parse()
	var err error
	if i2c, err = New(DS1307_I2C_ADDRESS, 2); err != nil {
		panic(err.Error())
	}
}

func decToBcd(val byte) byte {
	return ((val / 10 * 16) + (val % 10))
}

func bcdToDec(val byte) byte {
	return ((val / 16 * 10) + (val % 16))
}

func setDateDS1307(now *time.Time) error {
	// 0 to bit 7 starts the clock
	// If you want 12 hour am/pm you need to set bit 6 of (hour) (also need to change getDateDS1307)
	second, minute, hour, dayOfWeek, dayOfMonth, month, year := byte(now.Second()), byte(now.Minute()), byte(now.Hour()), byte(now.Weekday()), byte(now.Day()), byte(now.Month()), byte(now.Year()%2000)
	return i2c.Write(0,
		decToBcd(second),
		decToBcd(minute),
		decToBcd(hour),
		decToBcd(dayOfWeek),
		decToBcd(dayOfMonth),
		decToBcd(month),
		decToBcd(year))
}

func getDateDS1307() *time.Time {
	i2c.Write(0)
	b := make([]byte, 7)
	if err := i2c.Read(b); err != nil {
		fmt.Fprintf(os.Stderr, "getDateDS1307: %v\n", err)
		return nil
	}
	// A few of these need masks because certain bits are control bits
	second := bcdToDec(b[0] & 0x7f)
	minute := bcdToDec(b[1])
	hour := bcdToDec(b[2] & 0x3f) // Need to change this if 12 hour am/pm
	dayOfMonth := bcdToDec(b[4])
	month := bcdToDec(b[5])
	year := bcdToDec(b[6])
	_, zone := time.Now().Zone()
	t, err := time.Parse(time.RFC3339, fmt.Sprintf("20%.2d-%.2d-%.2dT%.2d:%.2d:%.2d%+.2d:00", year, month, dayOfMonth, hour, minute, second, zone/3600))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return &t
}

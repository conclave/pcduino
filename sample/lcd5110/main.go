// +build linux

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	. "github.com/conclave/pcduino/core"
	"github.com/conclave/pcduino/module/pcd8544"
)

func init() {
	Init()
	hostname, _ = os.Hostname()
	if hostname == "" {
		hostname = "pcDuino"
	}
	hostname += ":"
	fmt.Println("pcd8544+nokia5110 service")
	if conn, err := net.Dial("udp", "google.com:80"); err != nil {
		addr = "127.0.0.1"
	} else {
		addr = conn.LocalAddr().String()
		for i := 0; i < len(addr); i++ {
			if addr[i] == ':' {
				addr = addr[:i]
				break
			}
		}
	}
	fmt.Println(hostname, addr)
}

var lcd *pcd8544.LCD
var hostname string
var addr string
var step int64

func main() {
	flag.Int64Var(&step, "step", 5000, "clock step")
	flag.Parse()
	lcd = pcd8544.New(1, 0, 2, 4, 3, 60)
	lcd.Init()
	Delay(500)
	for {
		loop()
	}
}

func loop() {
	now := time.Now()
	Read()
	lcd.Clear()
	lcd.DrawString(0, 0, hostname)
	lcd.DrawLine(0, 9, 83, 9, pcd8544.BLACK)
	lcd.DrawString(0, 12, "UP "+singelton.Uptime.String())
	lcd.DrawString(0, 20, fmt.Sprintf("LD %2.1f %2.1f %2.1f", singelton.Loads[0], singelton.Loads[1], singelton.Loads[2]))
	lcd.DrawString(0, 28, fmt.Sprintf("%v %.2d:%.2d:%.2d", now.Weekday().String()[:3], now.Hour(), now.Minute(), now.Second()))
	lcd.DrawString(0, 36, addr)
	lcd.DrawLine(0, 45, 83, 45, pcd8544.BLACK)
	lcd.Display()
	Delay(step)
}

package main

import (
	. "github.com/conclave/pcduino/core"
	"github.com/conclave/pcduino/module/pcd8544"
)

func init() {
	Init()
}

func main() {
	lcd := pcd8544.New(1, 0, 2, 4, 3, 50)
	lcd.Init()
	lcd.Clear()
	lcd.ShowLogo()
	Delay(2000)
	for {
		loop()
	}
}

func loop() {
	Delay(200)
}

package main

import (
	"fmt"

	. "github.com/conclave/pcduino/core"
	. "github.com/conclave/pcduino/lib/i2c"
)

const (
	ADXL345_DEVID          = 0x00
	ADXL345_RESERVED1      = 0x01
	ADXL345_THRESH_TAP     = 0x1d
	ADXL345_OFSX           = 0x1e
	ADXL345_OFSY           = 0x1f
	ADXL345_OFSZ           = 0x20
	ADXL345_DUR            = 0x21
	ADXL345_LATENT         = 0x22
	ADXL345_WINDOW         = 0x23
	ADXL345_THRESH_ACT     = 0x24
	ADXL345_THRESH_INACT   = 0x25
	ADXL345_TIME_INACT     = 0x26
	ADXL345_ACT_INACT_CTL  = 0x27
	ADXL345_THRESH_FF      = 0x28
	ADXL345_TIME_FF        = 0x29
	ADXL345_TAP_AXES       = 0x2a
	ADXL345_ACT_TAP_STATUS = 0x2b
	ADXL345_BW_RATE        = 0x2c
	ADXL345_POWER_CTL      = 0x2d
	ADXL345_INT_ENABLE     = 0x2e
	ADXL345_INT_MAP        = 0x2f
	ADXL345_INT_SOURCE     = 0x30
	ADXL345_DATA_FORMAT    = 0x31
	ADXL345_DATAX0         = 0x32
	ADXL345_DATAX1         = 0x33
	ADXL345_DATAY0         = 0x34
	ADXL345_DATAY1         = 0x35
	ADXL345_DATAZ0         = 0x36
	ADXL345_DATAZ1         = 0x37
	ADXL345_FIFO_CTL       = 0x38
	ADXL345_FIFO_STATUS    = 0x39

	ADXL345_BW_1600 = 0xF // 1111
	ADXL345_BW_800  = 0xE // 1110
	ADXL345_BW_400  = 0xD // 1101
	ADXL345_BW_200  = 0xC // 1100
	ADXL345_BW_100  = 0xB // 1011
	ADXL345_BW_50   = 0xA // 1010
	ADXL345_BW_25   = 0x9 // 1001
	ADXL345_BW_12   = 0x8 // 1000
	ADXL345_BW_6    = 0x7 // 0111
	ADXL345_BW_3    = 0x6 // 0110

	/*
	   Interrupt PINs
	   INT1: 0
	   INT2: 1
	*/
	ADXL345_INT1_PIN = 0x00
	ADXL345_INT2_PIN = 0x01

	/*
	   Interrupt bit position
	*/
	ADXL345_INT_DATA_READY_BIT = 0x07
	ADXL345_INT_SINGLE_TAP_BIT = 0x06
	ADXL345_INT_DOUBLE_TAP_BIT = 0x05
	ADXL345_INT_ACTIVITY_BIT   = 0x04
	ADXL345_INT_INACTIVITY_BIT = 0x03
	ADXL345_INT_FREE_FALL_BIT  = 0x02
	ADXL345_INT_WATERMARK_BIT  = 0x01
	ADXL345_INT_OVERRUNY_BIT   = 0x00

	ADXL345_DATA_READY = 0x07
	ADXL345_SINGLE_TAP = 0x06
	ADXL345_DOUBLE_TAP = 0x05
	ADXL345_ACTIVITY   = 0x04
	ADXL345_INACTIVITY = 0x03
	ADXL345_FREE_FALL  = 0x02
	ADXL345_WATERMARK  = 0x01
	ADXL345_OVERRUNY   = 0x00

	ADXL345_SLAVE_ADDR = 0x53
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

var i2c *I2C

func setup() {
	var err error
	if i2c, err = New(ADXL345_SLAVE_ADDR, 2); err != nil {
		panic(err.Error())
	}
	fmt.Printf("dev id=0x%x\r\n", read8(0x0))
	//powerOn();
	write8(ADXL345_POWER_CTL, 0)
	write8(ADXL345_POWER_CTL, 16)
	write8(ADXL345_POWER_CTL, 8)

	//set activity/ inactivity thresholds (0-255)
	write8(ADXL345_THRESH_ACT, 75)   //setActivityThreshold(75); //62.5mg per increment
	write8(ADXL345_THRESH_INACT, 75) //setInactivityThreshold(75); //62.5mg per increment
	write8(ADXL345_TIME_INACT, 10)   //setTimeInactivity(10); // how many seconds of no activity is inactive?

	//look of activity movement on this axes - 1 == on; 0 == off
	setRegisterBit(ADXL345_ACT_INACT_CTL, 6, 1) //setActivityX(1);
	setRegisterBit(ADXL345_ACT_INACT_CTL, 5, 1) //setActivityY(1);
	setRegisterBit(ADXL345_ACT_INACT_CTL, 4, 1) //setActivityZ(1);

	//look of inactivity movement on this axes - 1 == on; 0 == off
	setRegisterBit(ADXL345_ACT_INACT_CTL, 2, 1) //setInactivityX(1);
	setRegisterBit(ADXL345_ACT_INACT_CTL, 1, 1) //setInactivityY(1);
	setRegisterBit(ADXL345_ACT_INACT_CTL, 0, 1) //setInactivityZ(1);

	//look of tap movement on this axes - 1 == on; 0 == off
	setRegisterBit(ADXL345_TAP_AXES, 2, 0) //setTapDetectionOnX(0);
	setRegisterBit(ADXL345_TAP_AXES, 1, 0) //setTapDetectionOnY(0);
	setRegisterBit(ADXL345_TAP_AXES, 0, 1) //setTapDetectionOnZ(1);

	//set values for what is a tap, and what is a double tap (0-255)
	write8(ADXL345_THRESH_TAP, 50) //setTapThreshold(50); //62.5mg per increment
	write8(ADXL345_DUR, 15)        //setTapDuration(15); //625¦Ìs per increment
	write8(ADXL345_LATENT, 80)     //setDoubleTapLatency(80); //1.25ms per increment
	write8(ADXL345_WINDOW, 200)    //setDoubleTapWindow(200); //1.25ms per increment

	//set values for what is considered freefall (0-255)
	write8(ADXL345_THRESH_FF, 7) //setFreeFallThreshold(7); //(5 - 9) recommended - 62.5mg per increment
	write8(ADXL345_TIME_FF, 45)  //setFreeFallDuration(45); //(20 - 70) recommended - 5ms per increment
}

func loop() {
	readXYZ()
	Delay(200000)
}

func write8(reg byte, value byte) {
	i2c.Write(reg, value)
}

func read8(reg byte) byte {
	b := []byte{0}
	i2c.Write(reg)
	i2c.Read(b)
	return b[0]
}

func read16(reg byte) uint16 {
	b := []byte{0, 0}
	i2c.Write(reg)
	i2c.Read(b)
	return uint16(b[0]) | uint16(b[1]<<8)
}

func readXYZ() {
	x := read16(ADXL345_DATAX0)
	y := read16(ADXL345_DATAY0)
	z := read16(ADXL345_DATAZ0)
	fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
}

func setRegisterBit(reg, bit, high byte) {
	value := read8(reg)
	if high != 0 {
		value |= (1 << bit)
	} else {
		value &= ^(1 << bit)
	}
	write8(reg, value)
}

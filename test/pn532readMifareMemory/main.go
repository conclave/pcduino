package main

import (
	"fmt"
	"os"

	. "github.com/conclave/pcduino/core"
	"github.com/conclave/pcduino/module/pn532"
)

const (
	SCK  = 13
	MOSI = 11
	SS   = 10
	MISO = 12
)

var nfc *pn532.PN532

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
	nfc = pn532.New(SCK, MISO, MOSI, SS)
	version := nfc.GetFirmwareVersion()
	if version == 0 {
		fmt.Fprintln(os.Stderr, "PN53x board not found")
		os.Exit(-1)
	}
	nfc.SAMConfig() // configure board to read RFID tags and cards
	fmt.Printf("Found chip PN5: %x\n", byte((version>>24)&0xFF))
	fmt.Printf("Firmware ver. %d.%d\n", byte((version>>16)&0xFF), byte((version>>8)&0xFF))
	fmt.Printf("Supports: %d\n", byte(version&0xFF))
}

func loop() {
	id := nfc.ReadPassiveTargetID(pn532.PN532_MIFARE_ISO14443A)
	if id != 0 {
		fmt.Printf("Read card #%d\n", id)
		keys := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		if nfc.AuthenticateBlock(1, byte(id), 0x08, pn532.KEY_A, keys) { //authenticate block 0x08
			block := make([]byte, 16)
			if nfc.ReadMemoryBlock(1, 0x08, block) {
				fmt.Printf("Read block_0x08: %x\n", block)
			}
		}
	}
	Delay(500)
}

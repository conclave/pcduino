package pn532

import (
	"bytes"
	"fmt"

	. "github.com/conclave/pcduino/core"
	. "github.com/conclave/pcduino/lib/spi"
)

const (
	PN532_PREAMBLE   = 0x00
	PN532_STARTCODE1 = 0x00
	PN532_STARTCODE2 = 0xFF
	PN532_POSTAMBLE  = 0x00

	PN532_HOSTTOPN532 = 0xD4

	PN532_FIRMWAREVERSION     = 0x02
	PN532_GETGENERALSTATUS    = 0x04
	PN532_SAMCONFIGURATION    = 0x14
	PN532_INLISTPASSIVETARGET = 0x4A
	PN532_INDATAEXCHANGE      = 0x40
	PN532_MIFARE_READ         = 0x30
	PN532_MIFARE_WRITE        = 0xA0

	PN532_AUTH_WITH_KEYA = 0x60
	PN532_AUTH_WITH_KEYB = 0x61

	PN532_WAKEUP = 0x55

	PN532_SPI_STATREAD  = 0x02
	PN532_SPI_DATAWRITE = 0x01
	PN532_SPI_DATAREAD  = 0x03
	PN532_SPI_READY     = 0x01

	PN532_MIFARE_ISO14443A = 0x0

	KEY_A = 1
	KEY_B = 2
)

type PN532 struct {
	clk    byte
	miso   byte
	mosi   byte
	ss     byte
	spi    *SPIDevice
	buffer []byte
}

var pn532ack = []byte{0x00, 0x00, 0xFF, 0x00, 0xFF, 0x00}
var pn532response_firmwarevers = []byte{0x00, 0xFF, 0x06, 0xFA, 0xD5, 0x03}

func New(clk, miso, mosi, ss byte) *PN532 {
	buffer := make([]byte, 64)
	spi := NewSPIDevice(0)
	spi.Begin()
	spi.SetDataMode(SPI_MODE0)
	spi.SetBitOrder(LSBFIRST)
	spi.SetClockDivider(SPI_CLOCK_DIV32)
	Delay(1000)
	p := &PN532{clk, miso, mosi, ss, spi, buffer}
	p.sendCommandCheckAck([]byte{PN532_FIRMWAREVERSION}, 1000) // ignore response
	return p
}

func (this *PN532) spiwrite(c byte) {
	this.spi.Transfer(c, SPI_CONTINUE)
}

func (this *PN532) spiwrite_end(c byte) {
	this.spi.Transfer(c, SPI_LAST)
}

func (this *PN532) spiread() byte {
	return this.spi.Transfer(0x00, SPI_CONTINUE)
}

func (this *PN532) spiread_end() byte {
	return this.spi.Transfer(0x00, SPI_LAST)
}

func (this *PN532) spiwritecommand(cmd []byte) {
	this.spiwrite(PN532_SPI_DATAWRITE)
	var checksum byte = PN532_PREAMBLE + PN532_PREAMBLE + PN532_STARTCODE2
	this.spiwrite(PN532_PREAMBLE)
	this.spiwrite(PN532_PREAMBLE)
	this.spiwrite(PN532_STARTCODE2)
	cmdlen := byte(len(cmd) + 1)
	this.spiwrite(cmdlen)
	this.spiwrite((^cmdlen) + 1)
	this.spiwrite(PN532_HOSTTOPN532)
	checksum += PN532_HOSTTOPN532
	for i := 0; i < len(cmd); i++ {
		this.spiwrite(cmd[i])
		checksum += cmd[i]
	}
	this.spiwrite(^checksum)
	this.spiwrite_end(PN532_POSTAMBLE)
}

func (this *PN532) readspidata(buff []byte) {
	this.spiwrite(PN532_SPI_DATAREAD)
	n := len(buff)
	for i := 0; i < n-1; i++ {
		buff[i] = this.spiread()
	}
	buff[n-1] = this.spiread_end()
}

func (this *PN532) readspistatus() byte {
	this.spiwrite(PN532_SPI_STATREAD)
	return this.spiread_end()
}

func (this *PN532) spi_readack() bool {
	ackbuf := make([]byte, 6)
	this.readspidata(ackbuf)
	return bytes.Equal(ackbuf, pn532ack)
}

func (this *PN532) sendCommandCheckAck(cmd []byte, timeout int64) bool {
	var timer int64 = 0
	this.spiwritecommand(cmd)
	for this.readspistatus() != PN532_SPI_READY {
		if timeout != 0 {
			timer += 10
			if timer > timeout {
				return false
			}
		}
		Delay(10)
	}
	if !this.spi_readack() {
		return false
	}
	timer = 0
	for this.readspistatus() != PN532_SPI_READY {
		if timeout != 0 {
			timer += 10
			if timer > timeout {
				return false
			}
		}
		Delay(10)
	}
	return true // ack'd command
}

func (this *PN532) ReadPassiveTargetID(cardbaudrate byte) uint32 {
	this.buffer[0] = PN532_INLISTPASSIVETARGET
	this.buffer[1] = 1
	this.buffer[2] = cardbaudrate
	if !this.sendCommandCheckAck(this.buffer[:3], 1000) {
		return 0x00
	}
	this.readspidata(this.buffer[:20])
	if this.buffer[7] != 1 {
		return 0x00
	}
	var val uint16 = (uint16(this.buffer[9]) << 8) | uint16(this.buffer[10])
	fmt.Printf("Sens Response: 0x%x\n", val)
	fmt.Printf("Sel Response: 0x%x\n", this.buffer[11])
	cid := uint32(0)
	for i := 0; i < int(this.buffer[12]); i++ {
		cid <<= 8
		cid |= uint32(this.buffer[13+i])
	}
	return cid
}

// cardno:  1~2
// addr:    0~63
// block:   len == 16
func (this *PN532) WriteMemoryBlock(cardno byte, addr byte, block []byte) bool {
	this.buffer[0] = PN532_INDATAEXCHANGE
	this.buffer[1] = cardno
	this.buffer[2] = PN532_MIFARE_WRITE
	this.buffer[3] = addr
	copy(this.buffer[4:20], block)
	if !this.sendCommandCheckAck(this.buffer[:20], 1000) {
		return false
	}
	this.readspidata(this.buffer[:8])
	if this.buffer[6] == 0x41 && this.buffer[7] == 0x00 {
		return true
	}
	return false
}

func (this *PN532) ReadMemoryBlock(cardno byte, addr byte, block []byte) bool {
	this.buffer[0] = PN532_INDATAEXCHANGE
	this.buffer[1] = cardno
	this.buffer[2] = PN532_MIFARE_READ
	this.buffer[3] = addr
	if !this.sendCommandCheckAck(this.buffer[:4], 1000) {
		return false
	}
	this.readspidata(this.buffer[:24])
	if this.buffer[6] == 0x41 && this.buffer[7] == 0x00 {
		copy(block, this.buffer[8:24])
		return true
	}
	return false
}

func (this *PN532) AuthenticateBlock(cardno byte, addr byte, cid uint32, authtype byte, keys []byte) bool {
	this.buffer[0] = PN532_INDATAEXCHANGE
	this.buffer[1] = cardno
	if authtype == KEY_A {
		this.buffer[2] = PN532_AUTH_WITH_KEYA
	} else {
		this.buffer[2] = PN532_AUTH_WITH_KEYB
	}
	this.buffer[3] = addr
	copy(this.buffer[4:10], keys[:6])
	this.buffer[10] = (byte(cid>>24) & 0xFF)
	this.buffer[11] = (byte(cid>>16) & 0xFF)
	this.buffer[12] = (byte(cid>>8) & 0xFF)
	this.buffer[13] = (byte(cid>>0) & 0xFF)
	if !this.sendCommandCheckAck(this.buffer[:14], 1000) {
		return false
	}
	this.readspidata(this.buffer[:8])
	if this.buffer[6] == 0x41 && this.buffer[7] == 0x00 {
		return true
	}
	return false
}

func (this *PN532) SAMConfig() bool {
	this.buffer[0] = PN532_SAMCONFIGURATION
	this.buffer[1] = 0x01 // normal mode;
	this.buffer[2] = 0x14 // timeout 50ms * 20 = 1 second
	this.buffer[3] = 0x01 // use IRQ pin!
	if !this.sendCommandCheckAck(this.buffer[:4], 1000) {
		return false
	}
	this.readspidata(this.buffer[:8])
	return this.buffer[5] == 0x15
}

func (this *PN532) GetFirmwareVersion() uint32 {
	this.buffer[0] = PN532_FIRMWAREVERSION
	if !this.sendCommandCheckAck(this.buffer[:1], 1000) {
		return 0
	}
	this.readspidata(this.buffer[:12])
	if !bytes.Equal(this.buffer[:6], pn532response_firmwarevers) {
		return 0
	}
	response := uint32(this.buffer[6]) << 8
	response |= uint32(this.buffer[7])
	response <<= 8
	response |= uint32(this.buffer[8])
	response <<= 8
	response |= uint32(this.buffer[9])
	return response
}

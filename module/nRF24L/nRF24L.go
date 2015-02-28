package nRF24L

import (
	"fmt"

	"github.com/conclave/pcduino/core"
)

var tx_address [TX_ADR_WIDTH]byte = [TX_ADR_WIDTH]byte{0x34, 0x43, 0x10, 0x10, 0x01}

func Init() {
	core.PinMode(CE, core.OUTPUT)
	core.PinMode(SCK, core.OUTPUT)
	core.PinMode(CSN, core.OUTPUT)
	core.PinMode(MOSI, core.OUTPUT)
	core.PinMode(MISO, core.INPUT)
	core.PinMode(IRQ, core.INPUT)
	// init io
	core.DigitalWrite(IRQ, 0)
	core.DigitalWrite(CE, 0)
	core.DigitalWrite(CSN, 1)

	status := spiRead(STATUS) // read the mode’s status register, the default value should be 'E'
	fmt.Printf("status = %X\n", status)
}

func rxMode() {
	core.DigitalWrite(CE, 0)
	spiWriteBuf(WRITE_REG+RX_ADDR_P0, tx_address[:]) // Use the same address on the RX device as the TX device
	spiRwReg(WRITE_REG+EN_AA, 0x01)                  // Enable Auto.Ack:Pipe0
	spiRwReg(WRITE_REG+EN_RXADDR, 0x01)              // Enable Pipe0
	spiRwReg(WRITE_REG+RF_CH, 40)                    // Select RF channel 40
	spiRwReg(WRITE_REG+RX_PW_P0, TX_PLOAD_WIDTH)     // Select same RX payload width as TX Payload width
	spiRwReg(WRITE_REG+RF_SETUP, 0x07)               // TX_PWR:0dBm, Datarate:2Mbps, LNA:HCURR
	spiRwReg(WRITE_REG+CONFIG, 0x0f)                 // Set PWR_UP bit, enable CRC(2 unsigned chars) & Prim:RX. RX_DR enabled..
	core.DigitalWrite(CE, 1)                         // Set CE pin high to enable RX device
	//  This device is now ready to receive one packet of 16 unsigned chars payload from a TX device sending to address
	//  '3443101001', with auto acknowledgment, retransmit count of 10, RF channel 40 and datarate = 2Mbps.
}

func txMode() {
	core.DigitalWrite(CE, 0)
	spiWriteBuf(WRITE_REG+TX_ADDR, tx_address[:])    // Writes TX_Address to nRF24L01
	spiWriteBuf(WRITE_REG+RX_ADDR_P0, tx_address[:]) // RX_Addr0 same as TX_Adr for Auto.Ack
	spiRwReg(WRITE_REG+EN_AA, 0x01)                  // Enable Auto.Ack:Pipe0
	spiRwReg(WRITE_REG+EN_RXADDR, 0x01)              // Enable Pipe0
	spiRwReg(WRITE_REG+SETUP_RETR, 0x1a)             // 500us + 86us, 10 retrans...
	spiRwReg(WRITE_REG+RF_CH, 40)                    // Select RF channel 40
	spiRwReg(WRITE_REG+RF_SETUP, 0x07)               // TX_PWR:0dBm, Datarate:2Mbps, LNA:HCURR
	spiRwReg(WRITE_REG+CONFIG, 0x0e)                 // Set PWR_UP bit, enable CRC(2 unsigned chars) & Prim:TX. MAX_RT & TX_DS enabled..
	// spiWriteBuf(WR_TX_PLOAD, tx_buf[:])
	core.DigitalWrite(CE, 1)
}

func SetMode(mode byte) bool {
	switch mode {
	case MODE_RX:
		rxMode()
	case MODE_TX:
		txMode()
	default:
		return false
	}
	return true
}

func Recv(buf []byte) {
	status := spiRead(STATUS) // read register STATUS's value
	if status&RX_DR != 0 {    // if receive data ready (TX_DS) interrupt
		spiReadBuf(RD_RX_PLOAD, buf) // read playload to rx_buf
		spiRwReg(FLUSH_RX, 0)        // clear RX_FIFO
	}
	spiRwReg(WRITE_REG+STATUS, status) // clear RX_DR or TX_DS or MAX_RT interrupt flag
}

func Send(buf []byte) {
	status := spiRead(STATUS) // read register STATUS's value
	if status&TX_DS != 0 {    // if receive data ready (TX_DS) interrupt
		spiRwReg(FLUSH_TX, 0)
		spiWriteBuf(WR_TX_PLOAD, buf) // write playload to TX_FIFO
	}
	if status&MAX_RT != 0 {
		// if receive data ready (MAX_RT) interrupt, this is retransmit than  SETUP_RETR
		spiRwReg(FLUSH_TX, 0)
		spiWriteBuf(WR_TX_PLOAD, buf) // disable standy-mode
	}
	spiRwReg(WRITE_REG+STATUS, status) // clear RX_DR or TX_DS or MAX_RT interrupt flag
}

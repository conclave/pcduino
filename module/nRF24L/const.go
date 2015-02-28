package nRF24L

const (
	MODE_RX = 0
	MODE_TX = 1
	//---------------------------------------------
	TX_ADR_WIDTH = 5
	// 5 unsigned chars TX(RX) address width
	TX_PLOAD_WIDTH = 32
	// 20 unsigned chars TX payload
	//---------------------------------------------
	CE = 8
	// CE_BIT:   Digital Input     Chip Enable Activates RX or TX mode
	CSN = 9
	// CSN BIT:  Digital Input     SPI Chip Select
	SCK = 10
	// SCK BIT:  Digital Input     SPI Clock
	MOSI = 11
	// MOSI BIT: Digital Input     SPI Slave Data Input
	MISO = 12
	// MISO BIT: Digital Output    SPI Slave Data Output, with tri-state option
	IRQ = 13
	// IRQ BIT:  Digital Output    Maskable interrupt pin

	//****************************************************
	// SPI(nRF24L01) commands
	READ_REG    = 0x00 // Define read command to register
	WRITE_REG   = 0x20 // Define write command to register
	RD_RX_PLOAD = 0x61 // Define RX payload register address
	WR_TX_PLOAD = 0xA0 // Define TX payload register address
	FLUSH_TX    = 0xE1 // Define flush TX register command
	FLUSH_RX    = 0xE2 // Define flush RX register command
	REUSE_TX_PL = 0xE3 // Define reuse TX payload register command
	NOP         = 0xFF // Define No Operation, might be used to read status register
	//***************************************************
	RX_DR  = 0x40
	TX_DS  = 0x20
	MAX_RT = 0x10
	//***************************************************
	// SPI(nRF24L01) registers(addresses)
	CONFIG      = 0x00 // 'Config' register address
	EN_AA       = 0x01 // 'Enable Auto Acknowledgment' register address
	EN_RXADDR   = 0x02 // 'Enabled RX addresses' register address
	SETUP_AW    = 0x03 // 'Setup address width' register address
	SETUP_RETR  = 0x04 // 'Setup Auto. Retrans' register address
	RF_CH       = 0x05 // 'RF channel' register address
	RF_SETUP    = 0x06 // 'RF setup' register address
	STATUS      = 0x07 // 'Status' register address
	OBSERVE_TX  = 0x08 // 'Observe TX' register address
	CD          = 0x09 // 'Carrier Detect' register address
	RX_ADDR_P0  = 0x0A // 'RX address pipe0' register address
	RX_ADDR_P1  = 0x0B // 'RX address pipe1' register address
	RX_ADDR_P2  = 0x0C // 'RX address pipe2' register address
	RX_ADDR_P3  = 0x0D // 'RX address pipe3' register address
	RX_ADDR_P4  = 0x0E // 'RX address pipe4' register address
	RX_ADDR_P5  = 0x0F // 'RX address pipe5' register address
	TX_ADDR     = 0x10 // 'TX address' register address
	RX_PW_P0    = 0x11 // 'RX payload width, pipe0' register address
	RX_PW_P1    = 0x12 // 'RX payload width, pipe1' register address
	RX_PW_P2    = 0x13 // 'RX payload width, pipe2' register address
	RX_PW_P3    = 0x14 // 'RX payload width, pipe3' register address
	RX_PW_P4    = 0x15 // 'RX payload width, pipe4' register address
	RX_PW_P5    = 0x16 // 'RX payload width, pipe5' register address
	FIFO_STATUS = 0x17 // 'FIFO Status Register' register address
)

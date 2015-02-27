// port of http://www.airspayce.com/mikem/arduino/VirtualWire/

package virtualwire

import (
	"fmt"

	. "github.com/conclave/pcduino/core"
)

/// By default the RX pin is expected to be low when idle, and to pulse high
/// for each data pulse.
/// This flag forces it to be inverted. This may be necessary if your transport medium
/// inverts the logic of your signal, such as happens with some types of A/V tramsmitter.
/// \param[in] inverted True to invert sense of receiver input
func (vw *VirtualWire) SetRxInverted(inverted bool) {
	vw.rx_inverted = inverted
}

/// By default the PTT pin goes high when the transmitter is enabled.
/// This flag forces it low when the transmitter is enabled.
/// \param[in] inverted True to invert PTT
func (vw *VirtualWire) SetPttInverted(inverted bool) {
	vw.ptt_inverted = inverted
	if vw.ptt_inverted {
		DigitalWrite(vw.ptt_pin, 1)
	} else {
		DigitalWrite(vw.ptt_pin, 0)
	}
}

/// Set vw to operate at speed bits per second
/// Must call StartRx() before you will get any messages
/// \param[in] speed Desired speed in bits per second
func (vw *VirtualWire) SetSpeed(speed uint16) {
}

func (vw *VirtualWire) StartTx() {
	vw.tx_index = 0
	vw.tx_bit = 0
	vw.tx_sample = 0

	// Enable the transmitter hardware
	if vw.ptt_inverted {
		DigitalWrite(vw.ptt_pin, 0)
	} else {
		DigitalWrite(vw.ptt_pin, 1)
	}

	// Next tick interrupt will send the first bit
	vw.tx_enabled = true
}

func (vw *VirtualWire) StopTx() {
	// Disable the transmitter hardware
	if vw.ptt_inverted {
		DigitalWrite(vw.ptt_pin, 1)
	} else {
		DigitalWrite(vw.ptt_pin, 0)
	}
	DigitalWrite(vw.tx_pin, 0)

	// No more ticks for the transmitter
	vw.tx_enabled = false
}

/// Start the Phase Locked Loop listening to the receiver
/// Must do this before you can receive any messages
/// When a message is available (good checksum or not), HaveMessage()
/// will return true.
func (vw *VirtualWire) StartRx() {
	if !vw.rx_enabled {
		vw.rx_enabled = true
		vw.rx_active = false // Never restart a partial message
	}
}

/// Stop the Phase Locked Loop listening to the receiver
/// No messages will be received until StartRx() is called again
/// Saves interrupt processing cycles
func (vw *VirtualWire) StopRx() {
	vw.rx_enabled = false
}

/// Returns the state of the
/// transmitter
/// \return true if the transmitter is active else false
func (vw *VirtualWire) ActiveTx() bool {
	return vw.tx_enabled
}

/// Block until the transmitter is idle
/// then returns
func (vw *VirtualWire) WaitTx() {
	for vw.tx_enabled {
	}
}

/// Block until a message is available
/// then returns
func (vw *VirtualWire) WaitRx() {
	for !vw.rx_done {
	}
}

/// Block until a message is available or for a max time
/// \param[in] milliseconds Maximum time to wait in milliseconds.
/// \return true if a message is available, false if the wait timed out.
func (vw *VirtualWire) WaitRxMax(milliseconds int64) bool {
	start := Millis()
	for !vw.rx_done && ((Millis() - start) < milliseconds) {
	}
	return vw.rx_done
}

/// Send a message with the given length. Returns almost immediately,
/// and message will be sent at the right timing by interrupts
/// \param[in] buf slice of the data to transmit
/// \return true if the message was accepted for transmission, false if the message is too long (>VW_MAX_MESSAGE_LEN - 3)
func (vw *VirtualWire) Send(buf []byte) error {
	l := len(buf)
	if l > VW_MAX_PAYLOAD {
		return fmt.Errorf("payload overflow")
	}
	var index int = VW_HEADER_LEN
	var crc uint16 = 0xffff
	var count = l + 3 // Added byte count and FCS to get total number of bytes

	// Wait for transmitter to become available
	vw.WaitTx()

	// Encode the message length
	crc = _crc_ccitt_update(crc, byte(count))
	vw.tx_buf[index] = symbols[count>>4]
	index++
	vw.tx_buf[index] = symbols[count&0xf]
	index++
	// Encode the message into 6 bit symbols. Each byte is converted into
	// 2 6-bit symbols, high nybble first, low nybble second
	for i := 0; i < l; i++ {
		crc = _crc_ccitt_update(crc, buf[i])
		vw.tx_buf[index] = symbols[buf[i]>>4]
		index++
		vw.tx_buf[index] = symbols[buf[i]&0xf]
		index++
	}

	// Append the fcs, 16 bits before encoding (4 6-bit symbols after encoding)
	// Caution: VW expects the _ones_complement_ of the CCITT CRC-16 as the FCS
	// VW sends FCS as low byte then hi byte
	crc = ^crc
	vw.tx_buf[index] = symbols[(crc>>4)&0xf]
	index++
	vw.tx_buf[index] = symbols[crc&0xf]
	index++
	vw.tx_buf[index] = symbols[(crc>>12)&0xf]
	index++
	vw.tx_buf[index] = symbols[(crc>>8)&0xf]
	index++

	// Total number of 6-bit symbols to send
	vw.tx_len = byte(index + VW_HEADER_LEN)

	// Start the low level interrupt handler sending symbols
	vw.StartTx()

	return nil
}

/// Returns true if an unread message is available
/// \return true if a message is available to read
func (vw *VirtualWire) HaveMessage() bool {
	return vw.rx_done
}

/// If a message is available (good checksum or not), copies
/// up to *len octets to buf.
func (vw *VirtualWire) GetMessage() ([]byte, error) {
	// Message available?
	if !vw.rx_done {
		return nil, fmt.Errorf("message unavailable")
	}

	// Wait until vw_rx_done is set before reading vw_rx_len
	// then remove bytecount and FCS
	var rxlen = vw.rx_len - 3
	buf := make([]byte, rxlen)
	copy(buf, vw.rx_buf[1:])

	vw.rx_done = false // OK, got that message thanks

	// Check the FCS, return goodness
	if vwCrc(vw.rx_buf[:vw.rx_len]) != 0xf0b8 { // FCS OK?
		return buf, fmt.Errorf("bad message")
	}
	return buf, nil
}

/// Returns the count of good messages received
/// Caution,: this is an 8 bit count and can easily overflow
/// \return Count of good messages received
func (vw *VirtualWire) GetRxGood() byte {
	return vw.rx_good
}

/// Returns the count of bad messages received, ie
/// messages with bogus lengths, indicating corruption
/// or lost octets.
/// Caution,: this is an 8 bit count and can easily overflow
/// \return Count of bad messages received
func (vw *VirtualWire) GetRxBad() byte {
	return vw.rx_bad
}

func NewVirtualWire(pins ...byte) *VirtualWire {
	var ptt, rx, tx byte = 10, 11, 12
	switch len(pins) {
	case 0:
	case 1:
		ptt = pins[0]
	case 2:
		ptt = pins[0]
		rx = pins[1]
	default:
		ptt = pins[0]
		rx = pins[1]
		tx = pins[2]
	}
	PinMode(tx, OUTPUT)
	PinMode(rx, INPUT)
	PinMode(ptt, OUTPUT)
	DigitalWrite(ptt, 0)
	return &VirtualWire{
		ptt_pin: ptt,
		rx_pin:  rx,
		tx_pin:  tx,
	}
}

type VirtualWire struct {
	ptt_pin        byte // The digital IO pin number of the press to talk, enables the transmitter hardware
	rx_pin         byte // The digital IO pin number of the receiver data
	tx_pin         byte // The digital IO pin number of the transmitter data
	tx_buf         [(VW_MAX_MESSAGE_LEN * 2) + VW_HEADER_LEN]byte
	tx_len         byte   // Number of symbols in vw_tx_buf to be sent;
	tx_index       byte   // Index of the next symbol to send. Ranges from 0 to tx_len
	tx_bit         byte   // Bit number of next bit to send
	tx_sample      byte   // Sample number for the transmitter. Runs 0 to 7 during one bit interval
	tx_enabled     bool   // Flag to indicated the transmitter is active
	tx_msg_count   uint16 // Total number of messages sent
	ptt_inverted   bool
	rx_inverted    bool
	rx_sample      byte // Current receiver sample
	rx_last_sample byte // Last receiver sample
	// PLL ramp, varies between 0 and VW_RX_RAMP_LEN-1 (159) over
	// VW_RX_SAMPLES_PER_BIT (8) samples per nominal bit time.
	// When the PLL is synchronised, bit transitions happen at about the
	// 0 mark.
	rx_pll_ramp   byte
	rx_integrator byte                     // This is the integrate and dump integral. If there are <5 0 samples in the PLL cycle the bit is declared a 0, else a 1
	rx_active     bool                     // Flag indicates if we have seen the start symbol of a new message and are in the processes of reading and decoding it
	rx_done       bool                     // Flag to indicate that a new message is available
	rx_enabled    bool                     // Flag to indicate the receiver PLL is to run
	rx_bits       uint16                   // Last 12 bits received, so we can look for the start symbol
	rx_bit_count  byte                     // How many bits of message we have received. Ranges from 0 to 12
	rx_buf        [VW_MAX_MESSAGE_LEN]byte // The incoming message buffer
	rx_count      byte                     // The incoming message expected length
	rx_len        byte                     // The incoming message buffer length received so far
	rx_bad        byte                     // Number of bad messages received and dropped due to bad lengths
	rx_good       byte                     // Number of good messages received
}

// 4 bit to 6 bit symbol converter table
// Used to convert the high and low nybbles of the transmitted data
// into 6 bit symbols for transmission. Each 6-bit symbol has 3 1s and 3 0s
// with at most 3 consecutive identical bits
var symbols []byte = []byte{0xd, 0xe, 0x13, 0x15, 0x16, 0x19, 0x1a, 0x1c,
	0x23, 0x25, 0x26, 0x29, 0x2a, 0x2c, 0x32, 0x34}

func vwSymbol6to4(symbol byte) byte {
	var count byte = 8

	// Linear search :-( Could have a 64 byte reverse lookup table?
	// There is a little speedup here courtesy Ralph Doncaster:
	// The shortcut works because bit 5 of the symbol is 1 for the last 8
	// symbols, and it is 0 for the first 8.
	// So we only have to search half the table
	for i := (symbol >> 2) & 8; count > 0; i++ {
		if symbol == symbols[i] {
			return i
		}
		count--
	}
	return 0 // Not found
}

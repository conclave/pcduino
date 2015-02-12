package virtualwire

const (
	/// Maximum number of bytes in a message, counting the byte count and FCS
	VW_MAX_MESSAGE_LEN = 80

	/// Number of samples per bit
	VW_RX_SAMPLES_PER_BIT = 8

	/// The maximum payload length
	VW_MAX_PAYLOAD = VW_MAX_MESSAGE_LEN - 3

	/// The size of the receiver ramp. Ramp wraps modulo this number
	VW_RX_RAMP_LEN = 160

	// Ramp adjustment parameters
	// Standard is if a transition occurs before VW_RAMP_TRANSITION (80) in the ramp,
	// the ramp is retarded by adding VW_RAMP_INC_RETARD (11)
	// else by adding VW_RAMP_INC_ADVANCE (29)
	// If there is no transition it is adjusted by VW_RAMP_INC (20)
	/// Internal ramp adjustment parameter
	VW_RAMP_INC = (VW_RX_RAMP_LEN / VW_RX_SAMPLES_PER_BIT)
	/// Internal ramp adjustment parameter
	VW_RAMP_TRANSITION = VW_RX_RAMP_LEN / 2
	/// Internal ramp adjustment parameter
	VW_RAMP_ADJUST = 9
	/// Internal ramp adjustment parameter
	VW_RAMP_INC_RETARD = (VW_RAMP_INC - VW_RAMP_ADJUST)
	/// Internal ramp adjustment parameter
	VW_RAMP_INC_ADVANCE = (VW_RAMP_INC + VW_RAMP_ADJUST)

	/// Outgoing message bits grouped as 6-bit words
	/// 36 alternating 1/0 bits, followed by 12 bits of start symbol
	/// Followed immediately by the 4-6 bit encoded byte count,
	/// message buffer and 2 byte FCS
	/// Each byte from the byte count on is translated into 2x6-bit words
	/// Caution, each symbol is transmitted LSBit first,
	/// but each byte is transmitted high nybble first
	VW_HEADER_LEN = 8
)

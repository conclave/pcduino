package virtualwire

func lo8(x uint16) byte { return (byte(x) & 0xff) }

func hi8(x uint16) byte { return (byte(x) >> 8) }

func crc16_update(crc uint16, a byte) uint16 {
	crc ^= uint16(a)
	for i := 0; i < 8; i++ {
		if crc&1 != 0 {
			crc = (crc >> 1) ^ 0xA001
		} else {
			crc = (crc >> 1)
		}
	}
	return crc
}

func crc_xmodem_update(crc uint16, data byte) uint16 {
	crc ^= (uint16(data) << 8)
	for i := 0; i < 8; i++ {
		if crc&0x8000 != 0 {
			crc = (crc << 1) ^ 0x1021
		} else {
			crc <<= 1
		}
	}
	return crc
}

func _crc_ccitt_update(crc uint16, data byte) uint16 {
	data ^= lo8(crc)
	data ^= data << 4
	return (((uint16(data) << 8) | uint16(hi8(crc))) ^ uint16(data>>4) ^ (uint16(data) << 3))
}

func _crc_ibutton_update(crc byte, data byte) byte {
	crc = crc ^ data
	for i := 0; i < 8; i++ {
		if crc&0x01 != 0 {
			crc = (crc >> 1) ^ 0x8C
		} else {
			crc >>= 1
		}
	}
	return crc
}

// Compute CRC over count bytes.
// This should only be ever called at user level, not interrupt level
func vwCrc(buf []byte) uint16 {
	crc := uint16(0xffff)
	count := len(buf)
	for i := 0; i < count; i++ {
		crc = _crc_ccitt_update(crc, buf[i])
	}
	return crc
}

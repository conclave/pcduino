package core

import (
	"unsafe"
)

// spidev.h
const SPI_IOC_MAGIC = 'k'

type SPI_ioc_transfer struct {
	tx_buf, rx_buf           uint64
	length, speed_hz         uint32
	delay_usecs              uint16
	bits_per_word, cs_change uint8
	padding                  uint32
}

// Read of SPI mode (SPI_MODE_0..SPI_MODE_3)
func SPI_IOC_RD_MODE() uintptr {
	return IOR(SPI_IOC_MAGIC, 1, 1)
}

// Write of SPI mode (SPI_MODE_0..SPI_MODE_3)
func SPI_IOC_WR_MODE() uintptr {
	return IOW(SPI_IOC_MAGIC, 1, 1)
}

// Read SPI bit justification
func SPI_IOC_RD_LSB_FIRST() uintptr {
	return IOR(SPI_IOC_MAGIC, 2, 1)
}

// Write SPI bit justification
func SPI_IOC_WR_LSB_FIRST() uintptr {
	return IOW(SPI_IOC_MAGIC, 2, 1)
}

// Read SPI device word length (1..N)
func SPI_IOC_RD_BITS_PER_WORD() uintptr {
	return IOR(SPI_IOC_MAGIC, 3, 1)
}

// Write SPI device word length (1..N)
func SPI_IOC_WR_BITS_PER_WORD() uintptr {
	return IOW(SPI_IOC_MAGIC, 3, 1)
}

// Read SPI device default max speed hz
func SPI_IOC_RD_MAX_SPEED_HZ() uintptr {
	return IOR(SPI_IOC_MAGIC, 4, 4)
}

// Write SPI device default max speed hz
func SPI_IOC_WR_MAX_SPEED_HZ() uintptr {
	return IOW(SPI_IOC_MAGIC, 4, 4)
}

// Write custom SPI message
func SPI_IOC_MESSAGE(n uintptr) uintptr {
	return IOW(SPI_IOC_MAGIC, 0, uintptr(SPI_MESSAGE_SIZE(n)))
}
func SPI_MESSAGE_SIZE(n uintptr) uintptr {
	if (n * unsafe.Sizeof(SPI_ioc_transfer{})) < (1 << IOC_SIZEBITS) {
		return (n * unsafe.Sizeof(SPI_ioc_transfer{}))
	}
	return 0
}

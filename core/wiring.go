package core

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	. "github.com/conclave/pcduino/sunxi"
)

const (
	GPIO_MODE_DIR           = "/sys/devices/virtual/misc/gpio/mode/"
	GPIO_PIN_DIR            = "/sys/devices/virtual/misc/gpio/pin/"
	GPIO_IF_PREFIX          = "gpio"
	ADC_IF                  = "/proc/adc"
	MAX_PWM_LEVEL           = 255
	EXTERNAL_NUM_INTERRUPTS = 2
)

const (
	PWMTMR_START    = 0x101
	PWMTMR_STOP     = 0x102
	PWMTMR_FUNC     = 0x103
	PWMTMR_TONE     = 0x104
	PWM_CONFIG      = 0x105
	HWPWM_DUTY      = 0x106
	PWM_FREQ        = 0x107
	MAX_PWMTMR_FREQ = 2000  //2kHz pin 3,9,10,11
	MIN_PWMTMR_FREQ = 126   //126Hz pin 3,9,10,11
	MAX_PWMHW_FREQ  = 20000 //20kHz pin 5,6
)

const pwm_dev = "/dev/pwmtimer"
const spi1_dev = "/dev/spidev1.0"
const spi2_dev = "/dev/spidev2.0"

var gpio_pin_fd []*os.File
var gpio_mode_fd []*os.File
var adc_fd []*os.File
var pwm_fd []*os.File

func init() {
	gpio_pin_fd = make([]*os.File, MAX_GPIO_NUM+1)
	gpio_mode_fd = make([]*os.File, MAX_GPIO_NUM+1)
	adc_fd = make([]*os.File, MAX_ADC_NUM+1)
}

// wiring

func Init() {
	var err error
	for i := 0; i <= MAX_GPIO_NUM; i++ {
		name := fmt.Sprintf("%s%s%d", GPIO_PIN_DIR, GPIO_IF_PREFIX, i)
		gpio_pin_fd[i], err = os.OpenFile(name, os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		name = fmt.Sprintf("%s%s%d", GPIO_MODE_DIR, GPIO_IF_PREFIX, i)
		gpio_mode_fd[i], err = os.OpenFile(name, os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}
	for i := 0; i <= 5; i++ { // why not MAX_ADC_NUM here?
		name := fmt.Sprintf("%s%d", ADC_IF, i)
		adc_fd[i], err = os.OpenFile(name, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}
}

// wiring_digital

func write_to_file(fd *os.File, data []byte) error {
	fd.Seek(0, os.SEEK_SET)
	n, err := fd.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}

func Hw_PinMode(pin, mode byte) {
	if pin >= 0 && pin <= MAX_GPIO_NUM && mode <= MAX_GPIO_MODE_NUM {
		data := []byte{mode + 0x30}
		if err := write_to_file(gpio_mode_fd[pin], data); err != nil {
			fmt.Fprintf(os.Stderr, "write gpio %d mode failed: %v\n", pin, err)
			os.Exit(-1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "hw_pinmode error: invalid pin or mode, pin=%d, mode=%d\n", pin, mode)
		os.Exit(-1)
	}
}

func PinMode(pin, mode byte) {
	switch pin {
	case 3, 9, 10, 11:
		fd, err := syscall.Open("/dev/pwmtimer", os.O_RDONLY|syscall.O_CLOEXEC, 0666)
		if err != nil {
			panic("open pwm device fail")
		}
		defer syscall.Close(fd)
		var val uint32 = uint32(pin)
		if err = Ioctl(fd, 0x102, uintptr(unsafe.Pointer(&val))); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			panic("can't set PWMTMR_STOP")
		}
	}
	switch mode {
	case INPUT, OUTPUT:
		Hw_PinMode(pin, mode)
	case INPUT_PULLUP:
		Hw_PinMode(pin, 8)
	}
}

func DigitalWrite(pin, value byte) {
	if pin >= 0 && pin <= MAX_GPIO_NUM && (value == HIGH || value == LOW) {
		data := []byte{value + 0x30}
		if err := write_to_file(gpio_pin_fd[pin], data); err != nil {
			fmt.Fprintf(os.Stderr, "write gpio %d failed: %v\n", pin, err)
			os.Exit(-1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "digitalWrite error: invalid pin or value, pin=%d, value=%d\n", pin, value)
		os.Exit(-1)
	}
}

func DigitalRead(pin byte) byte {
	if pin >= 0 && pin <= MAX_GPIO_NUM {
		fd := gpio_pin_fd[pin]
		b := make([]byte, 1)
		fd.Seek(0, os.SEEK_SET)
		if _, err := fd.Read(b); err != nil {
			fmt.Fprintf(os.Stderr, "read gpio %d failed\n", pin)
			os.Exit(-1)
		}
		b[0] -= '0'
		switch b[0] {
		case LOW:
			return LOW
		case HIGH:
			return HIGH
		default:
			return 0xFF
		}
	} else {
		fmt.Fprintf(os.Stderr, "digitalRead error: invalid pin, pin=%d\n", pin)
		os.Exit(-1)
	}
	return 0xFF
}

// wiring_pulse

//caution: if the pulse and timeout is too large, the CPU will continue 100% usage until the func quit.
func PulseIn(pin, state byte, timeout int64) int64 {
	PinMode(pin, INPUT)
	start := Micros()
	for DigitalRead(pin) == state {
		if Micros()-start > timeout {
			return 0
		}
	}
	for DigitalRead(pin) != state {
		if Micros()-start > timeout {
			return 0
		}
	}
	value := Micros()
	for DigitalRead(pin) == state {
		if Micros()-start > timeout {
			return 0
		}
	}
	return Micros() - value
}

// wiring_shift

func ShiftIn(dataPin, clockPin, bitOrder byte) byte {
	var value byte = 0
	PinMode(clockPin, OUTPUT)
	PinMode(dataPin, INPUT)
	for i := uint(0); i < 8; i++ {
		DigitalWrite(clockPin, HIGH)
		if bitOrder == LSBFIRST {
			value |= (DigitalRead(dataPin) << i)
		} else {
			value |= (DigitalRead(dataPin) << (7 - i))
		}
		DigitalWrite(clockPin, LOW)
	}
	return value
}

func ShiftOut(dataPin, clockPin, bitOrder, value byte) {
	PinMode(clockPin, OUTPUT)
	PinMode(dataPin, OUTPUT)
	for i := uint(0); i < 8; i++ {
		if bitOrder == LSBFIRST {
			DigitalWrite(dataPin, (value&(1<<i))&0x01)
		} else {
			DigitalWrite(dataPin, (value&(1<<(7-i)))&0x01)
		}
		DigitalWrite(clockPin, HIGH)
		DigitalWrite(clockPin, LOW)
	}
}

// wiring_analog

type PWM_Config struct {
	channel   int
	dutycycle int
}

type PWM_Freq struct {
	channel int
	step    int
	scale   int
	freq    uint
}

func SPI_adc_read_data(channel byte) int {
	PinMode(SPIEX_CS, OUTPUT)
	Hw_PinMode(SPIEX_MOSI, IO_SPIEX_FUNC)
	Hw_PinMode(SPIEX_MISO, IO_SPIEX_FUNC)
	Hw_PinMode(SPIEX_CLK, IO_SPIEX_FUNC)
	transfer := SPI_IOC_Transfer{}
	mode := 0
	fd, err := syscall.Open(spi2_dev, os.O_RDWR|syscall.O_CLOEXEC, 0666)
	if err != nil {
		fd, err = syscall.Open(spi1_dev, os.O_RDWR|syscall.O_CLOEXEC, 0666)
	}
	if err != nil {
		panic("can't open spi device: " + err.Error())
	}

	// MODE0
	if err = Ioctl(fd, SPI_IOC_RD_MODE(), uintptr(unsafe.Pointer(&mode))); err != nil {
		panic("can't get spi mode: " + err.Error())
	}
	mode &= ^0x3
	if err = Ioctl(fd, SPI_IOC_WR_MODE(), uintptr(unsafe.Pointer(&mode))); err != nil {
		panic("can't set spi mode: " + err.Error())
	}

	// MSBFIRST
	mode = 0
	if err = Ioctl(fd, SPI_IOC_WR_LSB_FIRST(), uintptr(unsafe.Pointer(&mode))); err != nil {
		panic("can't set bits order: " + err.Error())
	}
	DigitalWrite(SPIEX_CS, 0)
	var reg_val int = 0xC300 | (int(channel&0x7) << 11)
	txbuf := make([]byte, 2)
	rxbuf := make([]byte, 2)
	txbuf[0] = byte((reg_val >> 8)) & 0xFF
	txbuf[1] = byte(reg_val) & 0xFF
	transfer.TX_buf = uint64(uintptr(unsafe.Pointer(&txbuf)))
	transfer.Length = 2
	transfer.Speed_hz = 1000000
	transfer.Bits_per_word = 8
	transfer.Delay_usecs = 0xFFFF
	if err = Ioctl(fd, SPI_IOC_MESSAGE(1), uintptr(unsafe.Pointer(&transfer))); err != nil {
		panic("can't send spi message: " + err.Error())
	}
	DigitalWrite(SPIEX_CS, 1)
	DelayMicrosends(10)
	DigitalWrite(SPIEX_CS, 0)
	DelayMicrosends(180)
	DigitalWrite(SPIEX_CS, 1)
	DelayMicrosends(10)
	DigitalWrite(SPIEX_CS, 0)
	transfer.TX_buf = 0
	transfer.RX_buf = uint64(uintptr(unsafe.Pointer(&rxbuf)))
	transfer.Length = 2
	transfer.Speed_hz = 1000000
	transfer.Bits_per_word = 8
	transfer.Delay_usecs = 0xFFFF
	if err = Ioctl(fd, SPI_IOC_MESSAGE(1), uintptr(unsafe.Pointer(&transfer))); err != nil {
		panic("can't send spi message: " + err.Error())
	}
	DigitalWrite(SPIEX_CS, 1)
	syscall.Close(fd)
	return ((int(rxbuf[0]) << 8) | int(rxbuf[1])) >> 4
}

func AnalogRead(pin byte) int {
	if pin >= 0 && pin <= 5 {
		adc_fd[pin].Seek(0, os.SEEK_SET)
		b := make([]byte, 32)
		var s string
		if n, err := adc_fd[pin].Read(b); err != nil {
			fmt.Fprintf(os.Stderr, "read adc %d failed: %v\n", pin, err)
			os.Exit(-1)
		} else {
			s = string(b[:n])
		}
		// fmt.Printf("analogRead ret = %s\n", b)
		str := fmt.Sprintf("adc%d", pin)
		idx := strings.Index(s, str)
		if idx == -1 {
			return -1
		}
		idx += len(str) + 1
		s = s[idx : len(s)-1]
		if ret, err := strconv.Atoi(s); err != nil {
			return -1
		} else {
			return ret
		}
	} else if pin >= 6 && pin <= MAX_ADC_NUM {
		return SPI_adc_read_data(pin - 6)
	} else {
		fmt.Fprintf(os.Stderr, "analogRead error: invalid pin, pin=%d\n", pin)
		os.Exit(-1)
	}
	return -1
}

func PWMFreqSet(pin byte, freq uint) int {
	if (pin == 3 || pin == 5 || pin == 6 || pin == 9 || pin == 10 || pin == 11) && freq > 0 {
		pwmfreq := PWM_Freq{}
		pwmfreq.channel = int(pin)
		pwmfreq.freq = freq
		pwmfreq.step = 0
		fd, err := syscall.Open(pwm_dev, os.O_RDONLY|syscall.O_CLOEXEC, 0666)
		if err != nil {
			panic("open pwm device fail: " + err.Error())
		}
		switch pin {
		case 5, 6:
			if (freq == 195) || (freq == 260) || (freq == 390) || (freq == 520) || (freq == 781) {
				if err = Ioctl(fd, PWM_FREQ, uintptr(unsafe.Pointer(&pwmfreq))); err != nil {
					panic("can't set PWM_FREQ: " + err.Error())
				}
			} else {
				fmt.Fprintf(os.Stderr, "pwmfreq_set error: invalid frequency, should be [195,260,390,520,781], pin=%d\n", pin)
			}
		case 3, 9, 10, 11:
			if freq >= MIN_PWMTMR_FREQ && freq <= MAX_PWMTMR_FREQ {
				if err = Ioctl(fd, PWMTMR_STOP, uintptr(unsafe.Pointer(&pwmfreq.channel))); err != nil {
					panic("can't set PWMTMR_STOP: " + err.Error())
				}
				if err = Ioctl(fd, PWM_FREQ, uintptr(unsafe.Pointer(&pwmfreq))); err != nil {
					panic("can't set PWM_FREQ: " + err.Error())
				}
			} else {
				fmt.Fprintf(os.Stderr, "pwmfreq_set error: invalid frequency; pin=%d\n", pin)
			}
		}
		syscall.Close(fd)
		return pwmfreq.step
	} else {
		fmt.Fprintf(os.Stderr, "pwmfreq_set error: invalid pin, pin=%d\n", pin)
		os.Exit(-1)
	}
	return 0
}

func AnalogWrite(pin byte, value int) {
	if (pin == 3 || pin == 5 || pin == 6 || pin == 9 || pin == 10 || pin == 11) &&
		(value >= 0 && value <= MAX_PWM_LEVEL) {
		pwmconfig := PWM_Config{}
		pwmconfig.channel = int(pin)
		pwmconfig.dutycycle = value
		fd, err := syscall.Open(pwm_dev, os.O_RDONLY|syscall.O_CLOEXEC, 0666)
		if err != nil {
			panic("open pwm device fail: " + err.Error())
		}
		switch pin {
		case 5, 6:
			if err = Ioctl(fd, HWPWM_DUTY, uintptr(unsafe.Pointer(&pwmconfig))); err != nil {
				panic("can't set HWPWM_DUTY: " + err.Error())
			}
		case 3, 9, 10, 11:
			if err = Ioctl(fd, PWM_CONFIG, uintptr(unsafe.Pointer(&pwmconfig))); err != nil {
				panic("can't set PWM_CONFIG: " + err.Error())
			}
			var val int = 0
			if err = Ioctl(fd, PWMTMR_START, uintptr(unsafe.Pointer(&val))); err != nil {
				panic("can't set PWMTMR_START: " + err.Error())
			}
		}
		syscall.Close(fd)
	} else {
		fmt.Fprintf(os.Stderr, "analogWrite error: invalid pin, pin=%d\n", pin)
		os.Exit(-1)
	}
}

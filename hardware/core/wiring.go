package core

import (
	"fmt"
	"io"
	"os"
	"syscall"

	. "github.com/conclave/pcduino/hardware/sunxi"
)

const (
	GPIO_MODE_DIR           = "/sys/devices/virtual/misc/gpio/mode/"
	GPIO_PIN_DIR            = "/sys/devices/virtual/misc/gpio/pin/"
	GPIO_IF_PREFIX          = "gpio"
	ADC_IF                  = "/proc/adc"
	MAX_PWM_LEVEL           = 255
	EXTERNAL_NUM_INTERRUPTS = 2
)

// wiring_analog 1
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
	gpio_pin_fd = make([]*os.File, MAX_GPIO_NUM)
	gpio_mode_fd = make([]*os.File, MAX_GPIO_NUM)
	adc_fd = make([]*os.File, MAX_ADC_NUM)
}

// wiring
func Init() {
	var err error
	for i := 0; i < MAX_GPIO_NUM; i++ {
		name := fmt.Sprintf("%s%s%d", GPIO_PIN_DIR, GPIO_IF_PREFIX, i)
		gpio_pin_fd[i], err = os.OpenFile(name, os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s failed: %v\n", name, err)
			return
		}
		name = fmt.Sprintf("%s%s%d", GPIO_MODE_DIR, GPIO_IF_PREFIX, i)
		gpio_mode_fd[i], err = os.OpenFile(name, os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s failed: %v\n", name, err)
			return
		}
	}
	for i := 0; i < 5; i++ { // why not MAX_ADC_NUM here?
		name := fmt.Sprintf("%s%d", ADC_IF, i)
		adc_fd[i], err = os.OpenFile(name, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s failed: %v\n", name, err)
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

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

func Hw_PinMode(pin, mode byte) {
	if pin >= 0 && pin <= MAX_GPIO_NUM && mode <= MAX_GPIO_MODE_NUM {
		// data := []byte{mode}
		data := []byte(fmt.Sprintf("%d", mode))
		if err := write_to_file(gpio_mode_fd[pin], data); err != nil {
			fmt.Fprintf(os.Stderr, "write gpio %d mode failed\n", pin)
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
		fd, err := os.Open("/dev/pwmtimer")
		if err != nil {
			panic("open pwm device fail")
		}
		defer fd.Close()
		if err = ioctl(fd.Fd(), 0x102, uintptr(pin)); err != nil {
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
		data := []byte(fmt.Sprintf("%d", value))
		if err := write_to_file(gpio_pin_fd[pin], data); err != nil {
			fmt.Fprintf(os.Stderr, "write gpio %d failed\n", pin)
			os.Exit(-1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "digitalWrite error: invalid pin or mode, pin=%d, value=%d\n", pin, value)
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
		fmt.Fprintf(os.Stderr, "digitalRead error: invalid pin or mode, pin=%d\n", pin)
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
	for i := 0; i < 8; i++ {
		DigitalWrite(clockPin, HIGH)
		if bitOrder == LSBFIRST {
			value |= (DigitalRead(dataPin) << uint(i))
		} else {
			value |= (DigitalRead(dataPin) << uint(7-i))
		}
		DigitalWrite(clockPin, LOW)
	}
	return value
}

func ShiftOut(dataPin, clockPin, bitOrder, value byte) {
	PinMode(clockPin, OUTPUT)
	PinMode(dataPin, OUTPUT)
	for i := 0; i < 8; i++ {
		if bitOrder == LSBFIRST {
			DigitalWrite(dataPin, (value & (1 << uint(i))))
		} else {
			DigitalWrite(dataPin, (value & (1 << uint(7-i))))
		}
		DigitalWrite(clockPin, HIGH)
		DigitalWrite(clockPin, LOW)
	}
}

// wiring_analog 2

package lcd

import (
	. "github.com/conclave/pcduino/core"
)

const (
	// commands
	LCD_CLEARDISPLAY   = 0x01
	LCD_RETURNHOME     = 0x02
	LCD_ENTRYMODESET   = 0x04
	LCD_DISPLAYCONTROL = 0x08
	LCD_CURSORSHIFT    = 0x10
	LCD_FUNCTIONSET    = 0x20
	LCD_SETCGRAMADDR   = 0x40
	LCD_SETDDRAMADDR   = 0x80
	// flags for display entry mode
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00
	// flags for display on/off control
	LCD_DISPLAYON  = 0x04
	LCD_DISPLAYOFF = 0x00
	LCD_CURSORON   = 0x02
	LCD_CURSOROFF  = 0x00
	LCD_BLINKON    = 0x01
	LCD_BLINKOFF   = 0x00
	// flags for display/cursor shift
	LCD_DISPLAYMOVE = 0x08
	LCD_CURSORMOVE  = 0x00
	LCD_MOVERIGHT   = 0x04
	LCD_MOVELEFT    = 0x00
	// flags for function set
	LCD_8BITMODE = 0x10
	LCD_4BITMODE = 0x00
	LCD_2LINE    = 0x08
	LCD_1LINE    = 0x00
	LCD_5x10DOTS = 0x04
	LCD_5x8DOTS  = 0x00
)

type LCD struct {
	rs             byte
	rw             byte
	ep             byte
	date_pins      []byte
	fn, ctrl, mode byte
	lines, curln   byte
}

func New(rs, rw, byte, ep byte, ds ...byte) *LCD {
	l := len(ds)
	if l != 4 && l != 8 {
		return nil
	}
	lcd := LCD{
		rs, rw, ep, nil, 0, 0, 0, 0, 0,
	}
	PinMode(rs, OUTPUT)
	if rw != 0xFF {
		PinMode(rw, OUTPUT)
	}
	PinMode(ep, OUTPUT)
	if l == 4 {
		lcd.fn = LCD_4BITMODE | LCD_1LINE | LCD_5x8DOTS
	} else {
		lcd.fn = LCD_8BITMODE | LCD_1LINE | LCD_5x8DOTS
	}
	lcd.Begin(16, 1, LCD_5x8DOTS)
	return &lcd
}

func (this *LCD) Begin(cols, lines byte, dotsize byte) {
	if lines > 1 {
		this.fn |= LCD_2LINE
	}
	this.lines = lines
	this.curln = 0
	if dotsize != 0 && lines == 1 {
		this.fn |= LCD_5x10DOTS
	}
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way before 4.5V so we'll wait 50
	DelayMicrosends(50000)
	DigitalWrite(this.rs, LOW)
	DigitalWrite(this.ep, LOW)
	if this.rw != 0xFF {
		DigitalWrite(this.rw, LOW)
	}
	if this.fn&LCD_8BITMODE == 0 {
		this.write4bits(0x03)
		DelayMicrosends(4500)
		this.write4bits(0x03)
		DelayMicrosends(4500)
		this.write4bits(0x03)
		DelayMicrosends(150)
		this.write4bits(0x02)
	} else {
		this.command(LCD_FUNCTIONSET | this.fn)
		DelayMicrosends(4500)
		this.command(LCD_FUNCTIONSET | this.fn)
		DelayMicrosends(150)
		this.command(LCD_FUNCTIONSET | this.fn)
	}
	this.command(LCD_FUNCTIONSET | this.fn)
	this.ctrl = LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF
	this.Display()
	this.Clear()
	this.mode = LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT
	this.command(LCD_ENTRYMODESET | this.mode)
}

func (this *LCD) Clear() {
	this.command(LCD_CLEARDISPLAY)
	DelayMicrosends(2000)
}

func (this *LCD) Home() {
	this.command(LCD_RETURNHOME)
	DelayMicrosends(2000)
}

func (this *LCD) SetCursor(col, row byte) {
	row_offset := []byte{0x00, 0x40, 0x14, 0x54}
	if row > this.lines {
		row = this.lines - 1
	}
	this.command(LCD_SETDDRAMADDR | (col + row_offset[row]))
}

func (this *LCD) NoDisplay() {
	this.ctrl &= ^LCD_DISPLAYON + 0xFF + 1
	this.command(LCD_DISPLAYCONTROL | this.ctrl)
}

func (this *LCD) Display() {
	this.ctrl |= LCD_DISPLAYON
	this.command(LCD_DISPLAYCONTROL | this.ctrl)
}

func (this *LCD) NoCursor() {
	this.ctrl &= ^LCD_CURSORON + 0xFF + 1
	this.command(LCD_DISPLAYCONTROL | this.ctrl)
}

func (this *LCD) Cursor() {
	this.ctrl |= LCD_CURSORON
	this.command(LCD_DISPLAYCONTROL | this.ctrl)
}

func (this *LCD) NoBlink() {
	this.ctrl &= ^LCD_BLINKON + 0xFF + 1
	this.command(LCD_DISPLAYCONTROL | this.ctrl)
}

func (this *LCD) Blink() {
	this.ctrl |= LCD_BLINKON
	this.command(LCD_DISPLAYCONTROL | this.ctrl)
}

func (this *LCD) ScrollDisplayLeft() {
	this.command(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT)
}

func (this *LCD) ScrollDisplayRight() {
	this.command(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT)
}

// This is for text that flows Left to Right
func (this *LCD) LeftToRight() {
	this.mode |= LCD_ENTRYLEFT
	this.command(LCD_ENTRYMODESET | this.mode)
}

// This is for text that flows Right to Left
func (this *LCD) RightToLeft() {
	this.mode &= ^LCD_ENTRYLEFT + 0xFF + 1
	this.command(LCD_ENTRYMODESET | this.mode)
}

// This will 'right justify' text from the cursor
func (this *LCD) AutoScroll() {
	this.mode |= LCD_ENTRYSHIFTINCREMENT
	this.command(LCD_ENTRYMODESET | this.mode)
}

// This will 'left justify' text from the cursor
func (this *LCD) NoAutoScroll() {
	this.mode &= ^LCD_ENTRYSHIFTINCREMENT + 0xFF + 1
	this.command(LCD_ENTRYMODESET | this.mode)
}

func (this *LCD) CreateChar(location byte, charmap [8]byte) {
	location &= 0x7 // we only have 8 locations 0-7
	this.command(LCD_SETCGRAMADDR | (location << 3))
	for i := 0; i < 8; i++ {
		this.write(charmap[i])
	}
}

func (this *LCD) command(value byte) {
	this.send(value, LOW)
}

func (this *LCD) write(value byte) {
	this.send(value, HIGH)
}

func (this *LCD) send(value byte, mode byte) {
	DigitalWrite(this.rs, mode)
	if this.rw != 0xFF {
		DigitalWrite(this.rw, LOW)
	}
	if this.fn&LCD_8BITMODE != 0 {
		this.write8bits(value)
	} else {
		this.write4bits(value >> 4)
		this.write4bits(value)
	}
}

func (this *LCD) pulseEnable() {
	DigitalWrite(this.ep, LOW)
	DelayMicrosends(1)
	DigitalWrite(this.ep, HIGH)
	DelayMicrosends(1)
	DigitalWrite(this.ep, LOW)
	DelayMicrosends(100)
}

func (this *LCD) write4bits(value byte) {
	for i := byte(0); i < 4; i++ {
		PinMode(this.date_pins[i], OUTPUT)
		DigitalWrite(this.date_pins[i], (value>>i)&0x01)
	}
	this.pulseEnable()
}

func (this *LCD) write8bits(value byte) {
	for i := byte(0); i < 8; i++ {
		PinMode(this.date_pins[i], OUTPUT)
		DigitalWrite(this.date_pins[i], (value>>i)&0x01)
	}
	this.pulseEnable()
}

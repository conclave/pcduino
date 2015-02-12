package virtualwire

import (
	"testing"
)

func TestVirtualWirePins(t *testing.T) {
	vw := NewVirtualWire()
	if vw == nil {
		t.Fatal("cannot allocate!")
	}
	if vw.ptt_pin != 10 {
		t.Errorf("ptt pin not correct: %d\n", vw.ptt_pin)
	}
	if vw.rx_pin != 11 {
		t.Errorf("rx pin not correct: %d\n", vw.rx_pin)
	}
	if vw.tx_pin != 12 {
		t.Errorf("tx pin not correct: %d\n", vw.tx_pin)
	}
	vw = NewVirtualWire(5)
	if vw == nil {
		t.Fatal("cannot allocate!")
	}
	if vw.ptt_pin != 5 {
		t.Errorf("ptt pin not correct: %d\n", vw.ptt_pin)
	}
	if vw.rx_pin != 11 {
		t.Errorf("rx pin not correct: %d\n", vw.rx_pin)
	}
	if vw.tx_pin != 12 {
		t.Errorf("tx pin not correct: %d\n", vw.tx_pin)
	}
	vw = NewVirtualWire(5, 6)
	if vw == nil {
		t.Fatal("cannot allocate!")
	}
	if vw.ptt_pin != 5 {
		t.Errorf("ptt pin not correct: %d\n", vw.ptt_pin)
	}
	if vw.rx_pin != 6 {
		t.Errorf("rx pin not correct: %d\n", vw.rx_pin)
	}
	if vw.tx_pin != 12 {
		t.Errorf("tx pin not correct: %d\n", vw.tx_pin)
	}
	vw = NewVirtualWire(5, 6, 7)
	if vw == nil {
		t.Fatal("cannot allocate!")
	}
	if vw.ptt_pin != 5 {
		t.Errorf("ptt pin not correct: %d\n", vw.ptt_pin)
	}
	if vw.rx_pin != 6 {
		t.Errorf("rx pin not correct: %d\n", vw.rx_pin)
	}
	if vw.tx_pin != 7 {
		t.Errorf("tx pin not correct: %d\n", vw.tx_pin)
	}
}

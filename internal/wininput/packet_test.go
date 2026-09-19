package wininput

import (
	"encoding/binary"
	"testing"
)

func TestScanCodesForRequiredKeys(t *testing.T) {
	cases := []struct {
		key      uint32
		scan     uint16
		extended bool
	}{
		{0x5B, 0x5B, true},
		{0x5C, 0x5C, true},
		{0x08, 0x0E, false},
		{0xDC, 0x2B, false},
		{0xA2, 0x1D, false},
		{0xA3, 0x1D, true},
		{0x14, 0x3A, false},
	}
	for _, item := range cases {
		for _, down := range []bool{true, false} {
			packet, ok := Encode(item.key, down)
			if !ok {
				t.Fatalf("unsupported key %x", item.key)
			}
			if len(packet) != 40 {
				t.Fatal("Windows x64 INPUT must be 40 bytes")
			}
			if binary.LittleEndian.Uint32(packet[0:4]) != 1 {
				t.Fatal("INPUT type must be INPUT_KEYBOARD")
			}
			if binary.LittleEndian.Uint16(packet[8:10]) != 0 {
				t.Fatal("scan-code injection must not use a virtual key")
			}
			if binary.LittleEndian.Uint16(packet[10:12]) != item.scan {
				t.Fatalf("incorrect scan code for %x", item.key)
			}
			flags := binary.LittleEndian.Uint32(packet[12:16])
			if flags&8 == 0 || (flags&1 != 0) != item.extended || (flags&2 == 0) != down {
				t.Fatalf("incorrect flags %x for key %x down=%v", flags, item.key, down)
			}
			if binary.LittleEndian.Uint64(packet[24:32]) != Marker {
				t.Fatal("missing own-injection marker")
			}
		}
	}
}

func TestUnknownKeyIsRejected(t *testing.T) {
	if _, ok := Encode(0x43, true); ok {
		t.Fatal("unsupported key must not be injected")
	}
}

func TestOptionOutputsUseExplicitAltVirtualKeys(t *testing.T) {
	for _, key := range []uint32{0xA4, 0xA5} {
		for _, down := range []bool{true, false} {
			packet, ok := Encode(key, down)
			if !ok {
				t.Fatalf("Alt output %x was rejected", key)
			}
			if got := binary.LittleEndian.Uint16(packet[8:10]); got != uint16(key) {
				t.Fatalf("Alt virtual key = %x, want %x", got, key)
			}
			flags := binary.LittleEndian.Uint32(packet[12:16])
			if flags&0x0008 != 0 {
				t.Fatal("Alt output must use virtual-key injection")
			}
			if (flags&0x0001 != 0) != (key == 0xA5) || (flags&0x0002 == 0) != down {
				t.Fatalf("Alt side or release flags changed: %x", flags)
			}
		}
	}
}

package remap

import (
	"reflect"
	"testing"
)

// 0.10·HHKB ON에서 보고된 Mac 입력으로 추론한 변환(Ctrl → Command, Win → Control)을 검증한다.
func observedMacRole(key uint32) string {
	switch key {
	case 0xA2, 0xA3:
		return "Command"
	case 0xA4, 0xA5:
		return "Option"
	case 0x5B, 0x5C:
		return "Control"
	case 0x14:
		return "CapsLock"
	default:
		return "Other"
	}
}

func TestFinalMacModifierRolesWithObservedKeyboard(t *testing.T) {
	cases := []struct {
		source uint32
		hhkb   bool
		role   string
	}{
		{0xA4, true, "Command"},
		{0xA5, true, "Command"},
		{0x5B, true, "Option"},
		{0x5C, true, "Option"},
		{0x14, true, "Control"},
		{0xA2, true, "CapsLock"},
		{0xA3, true, "CapsLock"},
		{0xA4, false, "Command"},
		{0xA5, false, "Command"},
		{0x5B, false, "Option"},
		{0x5C, false, "Option"},
		{0x14, false, "CapsLock"},
		{0xA2, false, "Control"},
		{0xA3, false, "Control"},
	}
	for _, item := range cases {
		engine := New()
		engine.SetHHKBEnabled(item.hhkb)
		sink := &recorder{}
		for _, down := range []bool{true, false} {
			if !engine.Handle(item.source, down, false, sink.send) {
				sink.send(Output{Key: item.source, Down: down})
			}
		}
		if len(sink.events) != 2 || !sink.events[0].Down || sink.events[1].Down {
			t.Fatalf("modifier stroke mismatch: %v", sink.events)
		}
		for _, event := range sink.events {
			if role := observedMacRole(event.Key); role != item.role {
				t.Fatalf("source=%x hhkb=%t: Mac %s, want %s",
					item.source, item.hhkb, role, item.role)
			}
		}
		before := append([]Output(nil), sink.events...)
		for _, event := range before {
			if engine.Handle(event.Key, event.Down, true, sink.send) {
				t.Fatal("generated event entered the mapping cycle again")
			}
		}
		if !reflect.DeepEqual(sink.events, before) || engine.HasHeldInput() {
			t.Fatal("Parsec mapping generated recursive or held input")
		}
	}
}

package remap

import (
	"reflect"
	"testing"
)

func TestHHKBModeIsEnabledByDefault(t *testing.T) {
	engine := New()
	if !engine.HHKBEnabled() {
		t.Fatal("existing CapsLock/Control behavior must remain the default")
	}
}

func TestDisablingHHKBPassesCapsLockAndBothControlsThrough(t *testing.T) {
	engine := New()
	sink := &recorder{}
	if !engine.SetHHKBEnabled(false) {
		t.Fatal("could not disable idle HHKB mode")
	}
	for _, key := range []uint32{0x14, 0xA2, 0xA3} {
		target, mapped := engine.Target(key)
		if target != key || mapped {
			t.Fatalf("key %x must retain its original identity", key)
		}
		for _, down := range []bool{true, true, false} {
			if engine.Handle(key, down, false, sink.send) {
				t.Fatalf("original key %x was suppressed", key)
			}
		}
	}
	if len(sink.events) != 0 || engine.HasHeldInput() {
		t.Fatal("disabled HHKB mode generated input or left keys held")
	}
}

func TestOtherSwapsRemainEnabledWithoutHHKBMode(t *testing.T) {
	engine := New()
	engine.SetHHKBEnabled(false)
	sink := &recorder{}
	cases := []struct {
		source uint32
		target uint32
	}{
		{0xA4, 0x5B},
		{0xA5, 0x5C},
		{0x5B, 0xA4},
		{0x5C, 0xA5},
		{0xDC, 0x08},
		{0x08, 0xDC},
	}
	var want []Output
	for _, item := range cases {
		engine.Handle(item.source, true, false, sink.send)
		engine.Handle(item.source, false, false, sink.send)
		want = append(want,
			Output{Key: item.target, Down: true},
			Output{Key: item.target, Down: false},
		)
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("unrelated remaps changed: %v", sink.events)
	}
}

func TestHHKBModeCanBeEnabledAgain(t *testing.T) {
	engine := New()
	engine.SetHHKBEnabled(false)
	if !engine.SetHHKBEnabled(true) || !engine.HHKBEnabled() {
		t.Fatal("could not re-enable HHKB mode")
	}
	sink := &recorder{}
	engine.Handle(0x14, true, false, sink.send)
	engine.Handle(0x14, false, false, sink.send)
	want := []Output{
		{Key: 0xA2, Down: true},
		{Key: 0xA2, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("HHKB mapping did not recover: %v", sink.events)
	}
}

func TestHHKBChangeWaitsForHeldStrokeToFinish(t *testing.T) {
	for _, initiallyEnabled := range []bool{true, false} {
		engine := New()
		engine.SetHHKBEnabled(initiallyEnabled)
		sink := &recorder{}
		engine.Handle(0xA2, true, false, sink.send)
		if engine.SetHHKBEnabled(!initiallyEnabled) {
			t.Fatal("changed mapping halfway through a Control stroke")
		}
		if engine.HHKBEnabled() != initiallyEnabled {
			t.Fatal("rejected change still modified the setting")
		}
		blocked := engine.Handle(0xA2, false, false, sink.send)
		if blocked != initiallyEnabled {
			t.Fatal("key-up no longer matched key-down handling")
		}
		if !engine.SetHHKBEnabled(!initiallyEnabled) {
			t.Fatal("setting stayed blocked after key release")
		}
	}
}

func TestChangingHHKBOptionPreservesGlobalPause(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Pause(sink.send)
	engine.SetHHKBEnabled(false)
	if engine.Enabled() {
		t.Fatal("checkbox resumed a paused keyboard mapper")
	}
	engine.Resume()
	if engine.HHKBEnabled() {
		t.Fatal("resume reset the HHKB option")
	}
	if engine.Handle(0xA2, true, false, sink.send) {
		t.Fatal("Control was remapped after resuming with HHKB disabled")
	}
}

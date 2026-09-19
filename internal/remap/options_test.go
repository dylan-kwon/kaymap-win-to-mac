package remap

import (
	"reflect"
	"testing"
)

func TestHHKBModeIsEnabledByDefault(t *testing.T) {
	engine := New()
	if !engine.HHKBEnabled() {
		t.Fatal("existing HHKB swaps must remain the default")
	}
}

func TestDisablingHHKBPassesCapsAndBackspaceKeysThrough(t *testing.T) {
	engine := New()
	sink := &recorder{}
	if !engine.SetHHKBEnabled(false) {
		t.Fatal("could not disable idle HHKB mode")
	}
	for _, key := range []uint32{0x14, 0xDC, 0x08} {
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

func TestHHKBBackspaceSwapsRecoverWithRepeatedInput(t *testing.T) {
	for _, key := range []uint32{0xDC, 0x08} {
		engine := New()
		engine.SetHHKBEnabled(false)
		engine.SetHHKBEnabled(true)
		sink := &recorder{}
		target, _ := Target(key)
		for _, down := range []bool{true, true, false} {
			if !engine.Handle(key, down, false, sink.send) {
				t.Fatalf("HHKB key %x was not suppressed after re-enabling", key)
			}
		}
		want := []Output{
			{Key: target, Down: true},
			{Key: target, Down: true},
			{Key: target, Down: false},
		}
		if !reflect.DeepEqual(sink.events, want) || engine.HasHeldInput() {
			t.Fatalf("HHKB key %x repeat/release mismatch: %v", key, sink.events)
		}
	}
}

func TestHHKBChangeWaitsForHeldStrokeToFinish(t *testing.T) {
	for _, initiallyEnabled := range []bool{true, false} {
		for _, key := range []uint32{0x14, 0xA2, 0xA3, 0xDC, 0x08} {
			engine := New()
			engine.SetHHKBEnabled(initiallyEnabled)
			sink := &recorder{}
			engine.Handle(key, true, false, sink.send)
			if engine.SetHHKBEnabled(!initiallyEnabled) {
				t.Fatalf("changed mapping halfway through key %x stroke", key)
			}
			if engine.HHKBEnabled() != initiallyEnabled {
				t.Fatal("rejected change still modified the setting")
			}
			_, mapped := engine.Target(key)
			blocked := engine.Handle(key, false, false, sink.send)
			if blocked != mapped {
				t.Fatal("key-up no longer matched key-down handling")
			}
			if !engine.SetHHKBEnabled(!initiallyEnabled) {
				t.Fatal("setting stayed blocked after key release")
			}
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
	if !engine.Handle(0xA2, true, false, sink.send) {
		t.Fatal("Control must map to Ctrl output for Mac Control with HHKB disabled")
	}
	engine.Handle(0xA2, false, false, sink.send)
	want := []Output{
		{Key: 0xA2, Down: true},
		{Key: 0xA2, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("resumed Control mapping = %v", sink.events)
	}
}

func TestDisablingHHKBMapsPhysicalControlsToMacControl(t *testing.T) {
	for _, item := range []struct {
		source uint32
		target uint32
	}{
		{0xA2, 0xA2},
		{0xA3, 0xA3},
	} {
		engine := New()
		engine.SetHHKBEnabled(false)
		sink := &recorder{}
		for _, down := range []bool{true, false} {
			if !engine.Handle(item.source, down, false, sink.send) {
				t.Fatal("physical Control was not converted for Parsec swap mode")
			}
		}
		want := []Output{
			{Key: item.target, Down: true},
			{Key: item.target, Down: false},
		}
		if !reflect.DeepEqual(sink.events, want) || engine.HasHeldInput() {
			t.Fatalf("Control output = %v", sink.events)
		}
	}
}

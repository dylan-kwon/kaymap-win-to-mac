package remap

import (
	"reflect"
	"testing"
)

type recorder struct {
	events []Output
	fail   bool
}

func (r *recorder) send(event Output) bool {
	if r.fail {
		return false
	}
	r.events = append(r.events, event)
	return true
}

func TestRequiredKeySwapsPreserveDownAndUp(t *testing.T) {
	cases := []struct {
		name   string
		source uint32
		target uint32
	}{
		{"left alt", 0xA4, 0x5B},
		{"right alt", 0xA5, 0x5C},
		{"left windows", 0x5B, 0xA4},
		{"right windows", 0x5C, 0xA5},
		{"backslash", 0xDC, 0x08},
		{"backspace", 0x08, 0xDC},
		{"caps lock", 0x14, 0xA2},
		{"left control", 0xA2, 0x14},
		{"right control", 0xA3, 0x14},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			engine := New()
			sink := &recorder{}
			if !engine.Handle(item.source, true, false, sink.send) {
				t.Fatal("physical key-down was not suppressed")
			}
			if !engine.Handle(item.source, false, false, sink.send) {
				t.Fatal("physical key-up was not suppressed")
			}
			want := []Output{
				{Key: item.target, Down: true},
				{Key: item.target, Down: false},
			}
			if !reflect.DeepEqual(sink.events, want) {
				t.Fatalf("events = %v, want %v", sink.events, want)
			}
		})
	}
}

func TestShiftAndOrdinaryKeysPassThrough(t *testing.T) {
	engine := New()
	sink := &recorder{}
	for _, key := range []uint32{0xA0, 0xA1, 0x43, 0x09} {
		for _, down := range []bool{true, false} {
			if engine.Handle(key, down, false, sink.send) {
				t.Fatalf("unrelated key %x suppressed", key)
			}
		}
	}
	if len(sink.events) != 0 {
		t.Fatal("unrelated input generated synthetic events")
	}
}

func TestCapsLockActsAsHeldControlForShortcutAndPause(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0x14, true, false, sink.send)
	if engine.Handle(0x43, true, false, sink.send) {
		t.Fatal("C must pass through while mapped Control is held")
	}
	engine.Handle(0x43, false, false, sink.send)
	engine.Pause(sink.send)
	want := []Output{
		{Key: 0xA2, Down: true},
		{Key: 0xA2, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("Control chord events = %v", sink.events)
	}
	if !engine.Handle(0x14, false, false, sink.send) {
		t.Fatal("physical CapsLock release must be consumed")
	}
}

func TestBothControlsShareOneCapsLockPress(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0xA2, true, false, sink.send)
	engine.Handle(0xA2, true, false, sink.send)
	engine.Handle(0xA3, true, false, sink.send)
	engine.Handle(0xA2, false, false, sink.send)
	if len(sink.events) != 1 {
		t.Fatalf("repeat/overlap must not toggle or release CapsLock again: %v", sink.events)
	}
	engine.Handle(0xA3, false, false, sink.send)
	want := []Output{
		{Key: 0x14, Down: true},
		{Key: 0x14, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("shared CapsLock events = %v", sink.events)
	}
}

func TestPauseReleasesSharedCapsLockOnce(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0xA2, true, false, sink.send)
	engine.Handle(0xA3, true, false, sink.send)
	engine.Pause(sink.send)
	want := []Output{
		{Key: 0x14, Down: true},
		{Key: 0x14, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("shared CapsLock cleanup = %v", sink.events)
	}
	for _, key := range []uint32{0xA2, 0xA3} {
		if !engine.Handle(key, false, false, sink.send) {
			t.Fatal("orphan physical Control release escaped")
		}
	}
}

func TestRepeatedBackslashRepeatsDeletion(t *testing.T) {
	engine := New()
	sink := &recorder{}
	for range 3 {
		engine.Handle(0xDC, true, false, sink.send)
	}
	engine.Handle(0xDC, false, false, sink.send)
	want := []Output{
		{Key: 0x08, Down: true},
		{Key: 0x08, Down: true},
		{Key: 0x08, Down: true},
		{Key: 0x08, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("repeat events = %v", sink.events)
	}
}

func TestInjectedInputCannotTriggerReverseMapping(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0x14, true, false, sink.send)
	if engine.Handle(0xA2, true, true, sink.send) {
		t.Fatal("synthetic Ctrl key for Mac Control was remapped again")
	}
	if len(sink.events) != 1 {
		t.Fatal("recursive remapping generated extra events")
	}
}

func TestPauseReleasesHeldKeysAndConsumesTheirPhysicalRelease(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0x14, true, false, sink.send)
	engine.Handle(0xDC, true, false, sink.send)
	if !engine.Pause(sink.send) {
		t.Fatal("pause failed")
	}
	if engine.Enabled() {
		t.Fatal("mapping still enabled")
	}
	if len(sink.events) != 4 {
		t.Fatalf("held keys not released: %v", sink.events)
	}
	for _, key := range []uint32{0x14, 0xDC} {
		if !engine.Handle(key, true, false, sink.send) {
			t.Fatal("held key repeat escaped while paused")
		}
		if !engine.Handle(key, false, false, sink.send) {
			t.Fatal("orphan physical key-up escaped")
		}
	}
	if engine.Handle(0x14, true, false, sink.send) {
		t.Fatal("new input suppressed while paused")
	}
}

func TestResumeWaitsForPreviouslySuppressedKeys(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0x14, true, false, sink.send)
	engine.Pause(sink.send)
	engine.Resume()
	engine.Handle(0x14, true, false, sink.send)
	if len(sink.events) != 2 {
		t.Fatal("resume revived a held modifier")
	}
	engine.Handle(0x14, false, false, sink.send)
	engine.Handle(0x14, true, false, sink.send)
	if len(sink.events) != 3 {
		t.Fatal("fresh press was not mapped")
	}
}

func TestUnmatchedKeyUpPassesThrough(t *testing.T) {
	engine := New()
	sink := &recorder{}
	if engine.Handle(0x14, false, false, sink.send) {
		t.Fatal("key held before startup must be allowed to release")
	}
}

func TestInputReconnectMustWaitForPhysicalAndOutputRelease(t *testing.T) {
	engine := New()
	sink := &recorder{}
	if engine.HasHeldInput() {
		t.Fatal("new engine should be idle")
	}
	engine.Handle(0xA2, true, false, sink.send)
	if !engine.HasHeldInput() {
		t.Fatal("reconnect could interrupt a held mapped modifier")
	}
	engine.Pause(sink.send)
	if !engine.HasHeldInput() {
		t.Fatal("released output still has a physical key held")
	}
	engine.Handle(0xA2, false, false, sink.send)
	if engine.HasHeldInput() {
		t.Fatal("reconnect remained blocked after physical release")
	}
}

func TestFailedDownPassesOriginalStrokeThroughUntilRelease(t *testing.T) {
	engine := New()
	sink := &recorder{fail: true}
	if engine.Handle(0x14, true, false, sink.send) {
		t.Fatal("failed replacement swallowed original key-down")
	}
	sink.fail = false
	if engine.Handle(0x14, true, false, sink.send) {
		t.Fatal("partly remapped a stroke after initial send failure")
	}
	if engine.Handle(0x14, false, false, sink.send) {
		t.Fatal("original key-up swallowed after failed key-down")
	}
	if !engine.Handle(0x14, true, false, sink.send) {
		t.Fatal("next stroke did not recover")
	}
}

func TestFailedReleaseIsRetriedByPause(t *testing.T) {
	engine := New()
	sink := &recorder{}
	engine.Handle(0x14, true, false, sink.send)
	sink.fail = true
	engine.Handle(0x14, false, false, sink.send)
	if engine.Pause(sink.send) {
		t.Fatal("reported cleanup success despite send failure")
	}
	sink.fail = false
	if !engine.Pause(sink.send) {
		t.Fatal("release retry failed")
	}
	want := []Output{
		{Key: 0xA2, Down: true},
		{Key: 0xA2, Down: false},
	}
	if !reflect.DeepEqual(sink.events, want) {
		t.Fatalf("release retry events = %v", sink.events)
	}
}

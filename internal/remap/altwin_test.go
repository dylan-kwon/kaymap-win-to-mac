package remap

import (
	"reflect"
	"testing"
)

var modifierMappings = []struct {
	source uint32
	target uint32
}{
	{0xA4, 0x5B},
	{0xA5, 0x5C},
	{0x5B, 0xA4},
	{0x5C, 0xA5},
}

func TestFixedModifierMappingWithBothHHKBSettings(t *testing.T) {
	for _, hhkbEnabled := range []bool{true, false} {
		for _, item := range modifierMappings {
			engine := New()
			engine.SetHHKBEnabled(hhkbEnabled)
			target, mapped := engine.Target(item.source)
			if !mapped || target != item.target {
				t.Fatalf("source %x maps to %x, want %x", item.source, target, item.target)
			}
			sink := &recorder{}
			for _, key := range []uint32{0x41, 0x57} {
				if !engine.Handle(item.source, true, false, sink.send) {
					t.Fatal("physical modifier was not remapped")
				}
				for _, down := range []bool{true, false} {
					if engine.Handle(key, down, false, sink.send) {
						t.Fatal("A/W must pass through while mapped modifier is held")
					}
				}
				if !engine.Handle(item.source, false, false, sink.send) {
					t.Fatal("physical modifier release was not consumed")
				}
			}
			want := []Output{
				{Key: item.target, Down: true},
				{Key: item.target, Down: false},
				{Key: item.target, Down: true},
				{Key: item.target, Down: false},
			}
			if !reflect.DeepEqual(sink.events, want) || engine.HasHeldInput() {
				t.Fatalf("source %x shortcut events = %v", item.source, sink.events)
			}
		}
	}
}

func TestFixedModifierMappingSurvivesPauseAndResume(t *testing.T) {
	for _, item := range modifierMappings {
		engine := New()
		sink := &recorder{}
		engine.Handle(item.source, true, false, sink.send)
		engine.Pause(sink.send)
		engine.Resume()
		if !engine.Handle(item.source, false, false, sink.send) {
			t.Fatal("mapped modifier release escaped after pause")
		}
		engine.Handle(item.source, true, false, sink.send)
		engine.Handle(item.source, false, false, sink.send)
		want := []Output{
			{Key: item.target, Down: true},
			{Key: item.target, Down: false},
			{Key: item.target, Down: true},
			{Key: item.target, Down: false},
		}
		if !reflect.DeepEqual(sink.events, want) {
			t.Fatalf("source %x pause/resume events = %v", item.source, sink.events)
		}
		engine.Pause(sink.send)
		for _, down := range []bool{true, false} {
			if engine.Handle(item.source, down, false, sink.send) {
				t.Fatal("paused modifier must retain its original input")
			}
		}
		if !reflect.DeepEqual(sink.events, want) || engine.HasHeldInput() {
			t.Fatal("paused modifier generated synthetic input or left held keys")
		}
	}
}

func TestAltAndWindowsKeepSeparateModifierOutputs(t *testing.T) {
	for _, pair := range [][4]uint32{{0xA4, 0x5B, 0x5B, 0xA4}, {0xA5, 0x5C, 0x5C, 0xA5}} {
		for _, altFirst := range []bool{true, false} {
			for _, altReleasedFirst := range []bool{true, false} {
				engine := New()
				sink := &recorder{}
				presses := []uint32{pair[1], pair[0]}
				if altFirst {
					presses = []uint32{pair[0], pair[1]}
				}
				releases := []uint32{pair[1], pair[0]}
				if altReleasedFirst {
					releases = []uint32{pair[0], pair[1]}
				}
				var want []Output
				for _, down := range []bool{true, false} {
					keys := presses
					if !down {
						keys = releases
					}
					for _, key := range keys {
						target := pair[3]
						if key == pair[0] {
							target = pair[2]
						}
						if !engine.Handle(key, down, false, sink.send) {
							t.Fatalf("modifier %x escaped unchanged", key)
						}
						want = append(want, Output{Key: target, Down: down})
					}
				}
				if !reflect.DeepEqual(sink.events, want) || engine.HasHeldInput() {
					t.Fatalf("pair=%v altFirst=%t altReleasedFirst=%t: %v, want %v",
						pair, altFirst, altReleasedFirst, sink.events, want)
				}
			}
		}
	}
}

func TestPauseReleasesEveryMappedModifier(t *testing.T) {
	engine := New()
	sink := &recorder{}
	for _, item := range modifierMappings {
		engine.Handle(item.source, true, false, sink.send)
	}
	if !engine.Pause(sink.send) {
		t.Fatal("pause failed")
	}
	held := make(map[uint32]int)
	for _, event := range sink.events {
		if event.Down {
			held[event.Key]++
		} else {
			held[event.Key]--
		}
	}
	if len(sink.events) != 8 {
		t.Fatalf("pause did not release all four modifiers: %v", sink.events)
	}
	for key, count := range held {
		if count != 0 {
			t.Fatalf("output %x remained held", key)
		}
	}
	for _, item := range modifierMappings {
		if !engine.Handle(item.source, false, false, sink.send) {
			t.Fatal("physical modifier release escaped after cleanup")
		}
	}
	if engine.HasHeldInput() {
		t.Fatal("modifier state remained held after physical releases")
	}
}

func TestInjectedModifiersCannotTriggerReverseMapping(t *testing.T) {
	engine := New()
	sink := &recorder{}
	for _, item := range modifierMappings {
		for _, down := range []bool{true, false} {
			if engine.Handle(item.target, down, true, sink.send) {
				t.Fatal("injected modifier was remapped again")
			}
		}
	}
	if len(sink.events) != 0 || engine.HasHeldInput() {
		t.Fatal("injected modifiers produced recursive input")
	}
}

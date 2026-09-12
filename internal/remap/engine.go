// Package remap은 OS 입력 처리와 분리된 키 교환 상태를 관리한다.
package remap

type Output struct {
	Key  uint32
	Down bool
}

type Sender func(Output) bool

type Engine struct {
	enabled  bool
	held     map[uint32]uint32
	physical map[uint32]bool
	swallow  map[uint32]bool
	bypass   map[uint32]bool
}

func New() *Engine {
	return &Engine{
		enabled:  true,
		held:     make(map[uint32]uint32),
		physical: make(map[uint32]bool),
		swallow:  make(map[uint32]bool),
		bypass:   make(map[uint32]bool),
	}
}

func Target(key uint32) (uint32, bool) {
	switch key {
	case 0xA4:
		return 0x5B, true
	case 0x5B:
		return 0xA4, true
	case 0xA5:
		return 0x5C, true
	case 0x5C:
		return 0xA5, true
	case 0xDC:
		return 0x08, true
	case 0x08:
		return 0xDC, true
	case 0x14:
		return 0xA2, true
	case 0xA2, 0xA3:
		return 0x14, true
	default:
		return key, false
	}
}

func (e *Engine) otherSourceHolds(key uint32, target uint32) bool {
	for source, heldTarget := range e.held {
		if source != key && heldTarget == target {
			return true
		}
	}
	return false
}

// Handle은 원본 이벤트를 차단해야 하는 경우 true를 반환한다.
// Sender는 이벤트 삽입 성공 여부를 반환하며, 모든 호출은 같은 스레드에서 수행한다.
func (e *Engine) Handle(key uint32, down bool, injected bool, send Sender) bool {
	if injected {
		return false
	}
	target, mapped := Target(key)
	if !mapped {
		return false
	}
	if down {
		e.physical[key] = true
		if e.swallow[key] {
			return true
		}
		if e.bypass[key] {
			return false
		}
		if !e.enabled {
			e.bypass[key] = true
			return false
		}
		// 좌우 Ctrl이 같은 CapsLock을 공유하며, 반복 입력으로 다시 토글하지 않는다.
		_, alreadyHeld := e.held[key]
		if (target == 0x14 && alreadyHeld) || e.otherSourceHolds(key, target) {
			e.held[key] = target
			return true
		}
		if !send(Output{Key: target, Down: true}) {
			if _, held := e.held[key]; held {
				return true
			}
			e.bypass[key] = true
			return false
		}
		e.held[key] = target
		return true
	}
	delete(e.physical, key)
	if e.bypass[key] {
		delete(e.bypass, key)
		return false
	}
	wasSwallowed := e.swallow[key]
	delete(e.swallow, key)
	if heldTarget, held := e.held[key]; held {
		if e.otherSourceHolds(key, heldTarget) {
			delete(e.held, key)
			return true
		}
		if send(Output{Key: heldTarget, Down: false}) {
			delete(e.held, key)
		}
		return true
	}
	return wasSwallowed
}

// Pause는 눌러 둔 출력 키를 해제하고, 해당 물리 키가 떨어질 때까지 반복 입력을 차단한다.
// 해제 실패 항목은 다음 Pause 호출에서 재시도할 수 있도록 유지한다.
func (e *Engine) Pause(send Sender) bool {
	e.enabled = false
	success := true
	releases := make(map[uint32]bool)
	for key, target := range e.held {
		if e.physical[key] {
			e.swallow[key] = true
		}
		released, attempted := releases[target]
		if !attempted {
			released = send(Output{Key: target, Down: false})
			releases[target] = released
		}
		if released {
			delete(e.held, key)
		} else {
			success = false
		}
	}
	return success
}

func (e *Engine) Resume() {
	e.enabled = true
}

func (e *Engine) Enabled() bool {
	return e.enabled
}

func (e *Engine) HasHeldInput() bool {
	return len(e.physical) != 0 || len(e.held) != 0
}

package capture

import (
	"fmt"

	"kaymap/internal/remap"
)

type Diagnostic struct {
	count     uint64
	source    uint32
	target    uint32
	attempted bool
	sent      bool
}

// Record는 매핑 대상 키와 원본 Alt·Win 키의 마지막 누름만 메모리에 유지한다.
func (d *Diagnostic) Record(source uint32, target uint32, attempted bool, sent bool) {
	if _, mapped := remap.Target(source); !mapped {
		switch source {
		case 0xA4, 0xA5, 0x5B, 0x5C:
		default:
			return
		}
	}
	d.count++
	d.source = source
	d.target = target
	d.attempted = attempted
	d.sent = sent
}

func (d *Diagnostic) Text() string {
	if d.count == 0 {
		return "입력 대기 — 사용 중인 창에서 Win 키를 눌러 보세요."
	}
	state := "출력 전송 없음"
	if !d.attempted && d.source == d.target {
		state = "원본 입력 유지"
	}
	if d.attempted {
		state = "Windows 전송 실패"
		if d.sent {
			state = "Windows 전송 성공"
		}
	}
	return fmt.Sprintf("#%d  %s → %s / %s", d.count, keyName(d.source), keyName(d.target), state)
}

func keyName(key uint32) string {
	switch key {
	case 0x5B:
		return "Win(L)"
	case 0x5C:
		return "Win(R)"
	case 0xA4:
		return "Alt(L)"
	case 0xA5:
		return "Alt(R)"
	case 0xA2:
		return "Ctrl(L)"
	case 0xA3:
		return "Ctrl(R)"
	case 0x14:
		return "CapsLock"
	case 0x08:
		return "Backspace"
	case 0xDC:
		return "\\"
	default:
		return fmt.Sprintf("0x%X", key)
	}
}

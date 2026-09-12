// Package capture은 입력 연결 갱신 시점과 화면 진단 정보를 관리한다.
package capture

import "time"

type RefreshPlan struct {
	window   uintptr
	pending  bool
	due      time.Time
	manualAt time.Time
}

func (p *RefreshPlan) Request(at time.Time) {
	p.manualAt = at
	p.due = at
	p.pending = true
}

// Tick은 포커스 변경이 안정된 뒤, 키가 눌리지 않고 적용 중일 때 한 번 실행한다.
func (p *RefreshPlan) Tick(window uintptr, now time.Time, ready bool, attach func() bool) {
	if window == 0 {
		p.window = 0
		return
	}
	if window != p.window {
		p.window = window
		p.pending = true
		p.due = now.Add(500 * time.Millisecond)
		if p.manualAt.After(p.due) {
			p.due = p.manualAt
		}
	}
	if !p.pending || !ready || now.Before(p.due) {
		return
	}
	p.pending = false
	p.manualAt = time.Time{}
	attach()
}

package capture

import (
	"strings"
	"testing"
)

func TestDiagnosticShowsWindowsKeyAsAlt(t *testing.T) {
	var report Diagnostic
	report.Record(0x5B, 0xA4, true, true)
	text := report.Text()
	if !strings.Contains(text, "Win(L) → Alt(L)") || !strings.Contains(text, "Windows 전송 성공") {
		t.Fatalf("incorrect diagnosis: %s", text)
	}
	if strings.Contains(text, "Mac 전송 성공") {
		t.Fatal("local injection must not claim remote delivery")
	}
}

func TestDiagnosticDistinguishesFailedAndSuppressedInput(t *testing.T) {
	var report Diagnostic
	report.Record(0x5C, 0xA5, true, false)
	if !strings.Contains(report.Text(), "전송 실패") {
		t.Fatal("failed injection shown as success")
	}
	report.Record(0x5B, 0xA4, false, false)
	if !strings.Contains(report.Text(), "출력 전송 없음") {
		t.Fatal("paused or suppressed repeat shown as sent")
	}
}

func TestDiagnosticIgnoresOrdinaryKeys(t *testing.T) {
	var report Diagnostic
	report.Record(0x5B, 0xA4, true, true)
	before := report.Text()
	report.Record(0x41, 0x41, true, true)
	if report.Text() != before {
		t.Fatal("ordinary typing must not be recorded")
	}
}

func TestDiagnosticStartsWithNoInput(t *testing.T) {
	var report Diagnostic
	if report.Text() != "입력 대기 — 사용 중인 창에서 Win 키를 눌러 보세요." {
		t.Fatal("missing initial capture status")
	}
}

//go:build windows && (amd64 || arm64)

package main

import (
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"kaymap/internal/capture"
	"kaymap/internal/remap"
	"kaymap/internal/wininput"
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	gdi32             = syscall.NewLazyDLL("gdi32.dll")
	setHook           = user32.NewProc("SetWindowsHookExW")
	unhook            = user32.NewProc("UnhookWindowsHookEx")
	nextHook          = user32.NewProc("CallNextHookEx")
	sendInput         = user32.NewProc("SendInput")
	registerClass     = user32.NewProc("RegisterClassExW")
	createWindow      = user32.NewProc("CreateWindowExW")
	defaultWindowProc = user32.NewProc("DefWindowProcW")
	showWindow        = user32.NewProc("ShowWindow")
	getMessage        = user32.NewProc("GetMessageW")
	translateMessage  = user32.NewProc("TranslateMessage")
	dispatchMessage   = user32.NewProc("DispatchMessageW")
	isDialogMessage   = user32.NewProc("IsDialogMessageW")
	postMessage       = user32.NewProc("PostMessageW")
	postQuit          = user32.NewProc("PostQuitMessage")
	destroyWindow     = user32.NewProc("DestroyWindow")
	setText           = user32.NewProc("SetWindowTextW")
	sendMessage       = user32.NewProc("SendMessageW")
	getAsyncKeyState  = user32.NewProc("GetAsyncKeyState")
	getForeground     = user32.NewProc("GetForegroundWindow")
	setTimer          = user32.NewProc("SetTimer")
	killTimer         = user32.NewProc("KillTimer")
	messageBox        = user32.NewProc("MessageBoxW")
	loadCursor        = user32.NewProc("LoadCursorW")
	getModuleHandle   = kernel32.NewProc("GetModuleHandleW")
	createMutex       = kernel32.NewProc("CreateMutexW")
	closeHandle       = kernel32.NewProc("CloseHandle")
	getStockObject    = gdi32.NewProc("GetStockObject")
	app               application
)

const (
	wmClose         = 0x0010
	wmDestroy       = 0x0002
	wmCommand       = 0x0111
	wmInputFailure  = 0x8001
	pauseButtonID   = 101
	exitButtonID    = 102
	refreshButtonID = 103
	hhkbCheckboxID  = 104
	bmGetCheck      = 0x00F0
	bmSetCheck      = 0x00F1
	wmTimer         = 0x0113
)

type windowClass struct {
	Size        uint32
	Style       uint32
	Procedure   uintptr
	ClassExtra  int32
	WindowExtra int32
	Instance    uintptr
	Icon        uintptr
	Cursor      uintptr
	Background  uintptr
	MenuName    *uint16
	ClassName   *uint16
	SmallIcon   uintptr
}

type message struct {
	Window  uintptr
	ID      uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	X       int32
	Y       int32
	Private uint32
}

type keyboardEvent struct {
	Key   uint32
	Scan  uint32
	Flags uint32
	Time  uint32
	Extra uintptr
}

type application struct {
	engine          *remap.Engine
	window          uintptr
	status          uintptr
	pauseButton     uintptr
	hhkbCheckbox    uintptr
	hhkbInfo        uintptr
	hook            uintptr
	failed          bool
	instance        uintptr
	hookCallback    uintptr
	refresh         capture.RefreshPlan
	diagnostic      capture.Diagnostic
	diagnosticLabel uintptr
	diagnosticText  string
	attempted       bool
	sent            bool
}

func wide(text string) *uint16 {
	return syscall.StringToUTF16Ptr(text)
}

func setLabel(window uintptr, text string) {
	setText.Call(window, uintptr(unsafe.Pointer(wide(text))))
}

func alert(text string) {
	messageBox.Call(app.window, uintptr(unsafe.Pointer(wide(text))), uintptr(unsafe.Pointer(wide("Kaymap"))), 0x10)
}

func emit(event remap.Output) bool {
	app.attempted = true
	packet, ok := wininput.Encode(event.Key, event.Down)
	if !ok {
		return false
	}
	sent, _, _ := sendInput.Call(1, uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)))
	app.sent = sent == 1
	runtime.KeepAlive(packet)
	if sent != 1 && !app.failed {
		app.failed = true
		postMessage.Call(app.window, wmInputFailure, 0, 0)
	}
	return sent == 1
}

func keyboardHook(code int32, wParam uintptr, event *keyboardEvent) uintptr {
	if code == 0 {
		injected := event.Flags&0x10 != 0 || uint64(event.Extra) == wininput.Marker
		down := event.Flags&0x80 == 0
		if !injected {
			app.attempted = false
			app.sent = false
		}
		blocked := app.engine.Handle(event.Key, down, injected, emit)
		if !injected && down {
			target, _ := app.engine.Target(event.Key)
			app.diagnostic.Record(event.Key, target, app.attempted, app.sent)
		}
		if blocked {
			return 1
		}
	}
	result, _, _ := nextHook.Call(app.hook, uintptr(code), wParam, uintptr(unsafe.Pointer(event)))
	return result
}

func reconnectInput() bool {
	if !app.engine.Pause(emit) {
		setLabel(app.status, "입력 연결 실패 — 눌린 출력 키 해제 실패")
		setLabel(app.pauseButton, "다시 적용")
		return false
	}
	newHook, _, _ := setHook.Call(13, app.hookCallback, app.instance, 0)
	if newHook == 0 {
		setLabel(app.status, "입력 연결 실패 — 일시정지됨")
		setLabel(app.pauseButton, "다시 적용")
		return false
	}
	if app.hook != 0 {
		removed, _, _ := unhook.Call(app.hook)
		if removed == 0 {
			// Windows가 시간 초과로 기존 Hook을 제거했을 수도 있으므로 새 연결 유지.
			setLabel(app.status, "기존 입력 연결 확인 실패 — 새 연결로 적용 중")
		}
	}
	app.hook = newHook
	app.engine.Resume()
	setLabel(app.pauseButton, "일시정지")
	setLabel(app.status, "적용 중 — 입력 연결 갱신 완료")
	return true
}

func tickInput() {
	window, _, _ := getForeground.Call()
	ready := app.engine.Enabled() && !app.engine.HasHeldInput() && !mappedKeyIsDown()
	app.refresh.Tick(window, time.Now(), ready, reconnectInput)
	text := app.diagnostic.Text()
	if text != app.diagnosticText {
		setLabel(app.diagnosticLabel, text)
		app.diagnosticText = text
	}
}

func mappedKeyIsDown() bool {
	for _, key := range []uintptr{0xA4, 0xA5, 0x5B, 0x5C, 0x08, 0xDC, 0x14, 0xA2, 0xA3} {
		state, _, _ := getAsyncKeyState.Call(key)
		if state&0x8000 != 0 {
			return true
		}
	}
	return false
}

func pauseMapping() bool {
	clean := app.engine.Pause(emit)
	setLabel(app.pauseButton, "다시 적용")
	if clean {
		setLabel(app.status, "일시정지 — 원래 Windows 키 배열 사용 중")
	} else {
		setLabel(app.status, "키 해제 실패 — 모든 키를 떼고 다시 시도하세요.")
	}
	return clean
}

func toggleMapping() {
	if app.engine.Enabled() {
		pauseMapping()
		return
	}
	if !app.engine.Pause(emit) || mappedKeyIsDown() {
		setLabel(app.status, "매핑 대상 키를 모두 떼어 주세요.")
		return
	}
	app.failed = false
	app.engine.Resume()
	app.refresh.Request(time.Now().Add(500 * time.Millisecond))
	setLabel(app.status, "적용 중 — Windows 전체 키보드에 적용")
	setLabel(app.pauseButton, "일시정지")
}

func updateHHKBControls() {
	checked := uintptr(0)
	text := "Mac 기준: Ctrl → Control / CapsLock → CapsLock\r\nBackspace·역슬래시(\\) 원본 입력 유지"
	if app.engine.HHKBEnabled() {
		checked = 1
		text = "Mac 기준: CapsLock → Control / Ctrl → CapsLock\r\nBackspace ↔ 역슬래시(\\) 교환"
	}
	sendMessage.Call(app.hhkbCheckbox, bmSetCheck, checked, 0)
	setLabel(app.hhkbInfo, text)
}

func toggleHHKB() {
	checked, _, _ := sendMessage.Call(app.hhkbCheckbox, bmGetCheck, 0, 0)
	if mappedKeyIsDown() || !app.engine.SetHHKBEnabled(checked == 1) {
		updateHHKBControls()
		setLabel(app.status, "키를 모두 뗀 뒤 HHKB 옵션을 변경하세요.")
		return
	}
	updateHHKBControls()
	app.diagnostic = capture.Diagnostic{}
	app.diagnosticText = ""
	if app.engine.Enabled() {
		setLabel(app.status, "적용 중 — HHKB 옵션 변경 완료")
	} else {
		setLabel(app.status, "일시정지 — HHKB 옵션 변경 완료")
	}
}

func stopMapping() {
	if !pauseMapping() {
		alert("출력 키 해제에 실패했습니다. 모든 키를 떼고 종료를 다시 눌러 주세요.")
		return
	}
	if app.hook != 0 {
		unhook.Call(app.hook)
		app.hook = 0
	}
	killTimer.Call(app.window, 1)
	destroyWindow.Call(app.window)
}

func windowProcedure(window uintptr, id uint32, wParam uintptr, lParam uintptr) uintptr {
	switch id {
	case wmCommand:
		switch wParam & 0xFFFF {
		case pauseButtonID:
			toggleMapping()
			return 0
		case exitButtonID:
			stopMapping()
			return 0
		case refreshButtonID:
			app.refresh.Request(time.Now().Add(2 * time.Second))
			setLabel(app.status, "입력 연결 예약 — 적용 중인 상태로 사용할 창을 클릭하세요.")
			return 0
		case hhkbCheckboxID:
			if (wParam>>16)&0xFFFF == 0 {
				toggleHHKB()
			}
			return 0
		}
	case wmTimer:
		tickInput()
		return 0
	case wmInputFailure:
		pauseMapping()
		setLabel(app.status, "입력 전송 실패로 일시정지 — 대상 앱의 권한 확인 필요")
		return 0
	case wmClose:
		stopMapping()
		return 0
	case wmDestroy:
		postQuit.Call(0)
		return 0
	}
	result, _, _ := defaultWindowProc.Call(window, uintptr(id), wParam, lParam)
	return result
}

func addControl(class string, text string, x uintptr, y uintptr, width uintptr, height uintptr, id uintptr) uintptr {
	style := uintptr(0x50000000)
	if class == "BUTTON" {
		style |= 0x00010000
		if id == hhkbCheckboxID {
			style |= 0x00000003
		}
	}
	control, _, _ := createWindow.Call(
		0,
		uintptr(unsafe.Pointer(wide(class))),
		uintptr(unsafe.Pointer(wide(text))),
		style,
		x, y, width, height,
		app.window, id, 0, 0,
	)
	font, _, _ := getStockObject.Call(17)
	sendMessage.Call(control, 0x0030, font, 1)
	return control
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	mutex, _, mutexError := createMutex.Call(0, 0, uintptr(unsafe.Pointer(wide("Local\\KaymapPortable"))))
	if mutex == 0 {
		alert("프로그램 실행 상태를 확인할 수 없습니다.")
		return
	}
	defer closeHandle.Call(mutex)
	if mutexError == syscall.Errno(183) {
		alert("Kaymap이 이미 실행 중입니다. 기존 창을 사용하세요.")
		return
	}
	app.engine = remap.New()
	instance, _, _ := getModuleHandle.Call(0)
	app.instance = instance
	app.hookCallback = syscall.NewCallback(keyboardHook)
	cursor, _, _ := loadCursor.Call(0, 32512)
	class := windowClass{
		Procedure:  syscall.NewCallback(windowProcedure),
		Instance:   instance,
		Cursor:     cursor,
		Background: 16,
		ClassName:  wide("KaymapPortableWindow"),
	}
	class.Size = uint32(unsafe.Sizeof(class))
	registered, _, _ := registerClass.Call(uintptr(unsafe.Pointer(&class)))
	if registered == 0 {
		alert("사용 화면 생성에 실패했습니다.")
		return
	}
	app.window, _, _ = createWindow.Call(
		0,
		uintptr(unsafe.Pointer(class.ClassName)),
		uintptr(unsafe.Pointer(wide("Kaymap 1.0.0 — Parsec 매핑 ON"))),
		0x00CA0000,
		0x80000000, 0x80000000, 660, 450,
		0, 0, instance, 0,
	)
	if app.window == 0 {
		alert("사용 화면 생성에 실패했습니다.")
		return
	}
	app.status = addControl("STATIC", "적용 중 — Windows 전체 키보드에 적용", 20, 20, 610, 30, 0)
	addControl("STATIC", "현재 키보드 실측 기준\r\nAlt → Command / Win → Option\r\n실행 중 Windows 전체 키보드에 적용\r\n종료 버튼 또는 창 닫기로 키 매핑 해제", 20, 60, 610, 95, 0)
	app.hhkbCheckbox = addControl("BUTTON", "HHKB 모드", 20, 160, 610, 28, hhkbCheckboxID)
	app.hhkbInfo = addControl("STATIC", "", 40, 192, 590, 55, 0)
	updateHHKBControls()
	app.diagnosticLabel = addControl("STATIC", app.diagnostic.Text(), 20, 260, 610, 30, 0)
	addControl("STATIC", "Windows 전송 성공은 대상 앱의 수신 확인과 별개입니다.", 20, 295, 610, 25, 0)
	app.pauseButton = addControl("BUTTON", "일시정지", 20, 340, 150, 35, pauseButtonID)
	addControl("BUTTON", "입력 다시 연결", 185, 340, 180, 35, refreshButtonID)
	addControl("BUTTON", "종료", 380, 340, 150, 35, exitButtonID)
	if mappedKeyIsDown() {
		pauseMapping()
	}
	app.hook, _, _ = setHook.Call(13, app.hookCallback, instance, 0)
	if app.hook == 0 {
		alert("키보드 Hook 설치에 실패했습니다. 이 PC의 실행 제한을 확인하세요.")
		return
	}
	defer func() {
		if app.hook != 0 {
			unhook.Call(app.hook)
		}
	}()
	showWindow.Call(app.window, 1)
	timer, _, _ := setTimer.Call(app.window, 1, 100, 0)
	if timer == 0 {
		setLabel(app.status, "입력 자동 재연결 실패 — Kaymap 종료 후 다시 실행하세요.")
	}
	var event message
	for {
		result, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&event)), 0, 0, 0)
		if int32(result) <= 0 {
			break
		}
		handled, _, _ := isDialogMessage.Call(app.window, uintptr(unsafe.Pointer(&event)))
		if handled != 0 {
			continue
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&event)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&event)))
	}
	app.engine.Pause(emit)
}

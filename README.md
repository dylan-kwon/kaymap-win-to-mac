# Kaymap — 피씨방 Windows에서 Parsec으로 Mac 사용

- 설치·재부팅 없이 실행하는 Windows 10/11 x64용 포터블 키 매핑 프로그램
- 실행 파일: [kaymap-windows-x64.exe](dist/kaymap-windows-x64.exe)
- 배포 묶음: [kaymap-windows-x64.zip](dist/kaymap-windows-x64.zip), 실행 파일·사용법·Go 라이선스 포함
- 적용 범위: 실행 중 Windows 전체 키보드
- 종료 버튼 또는 창 닫기로 해제, 일시정지 및 다시 적용 지원
- Registry 변경·드라이버 설치·자동 시작 등록 없음
- 별도 Go·Python·AutoHotkey 설치 불필요
- 현재 단계: 키 매핑 로직 테스트 및 Windows 실행 파일 빌드 완료, Windows UI와 Parsec 실제 입력은 미검증

## 사용 절차

1. Windows PC로 [실행 파일](dist/kaymap-windows-x64.exe) 복사
2. Parsec 실행 및 Mac 연결
3. Windows Parsec → Settings → Client → `Swap Command and Ctrl for MacOS`를 `Off`로 설정
4. Parsec 세션에서 Keyboard Immersive Mode 활성화, 기본 토글 단축키 `Ctrl+Shift+I`
5. 모든 키를 뗀 상태에서 `kaymap-windows-x64.exe` 실행
6. 프로그램의 적용 상태 확인 후 Parsec 창 클릭
7. 사용 종료 시 Kaymap의 `종료` 버튼 클릭 또는 Kaymap 창 닫기

- 다운로드한 EXE와 함께 설치할 파일 없음
- 피씨방에서 파일이 초기화되는 경우 다음 방문 시 실행 파일만 다시 복사하여 실행
- SharpKeys 방식 제외: 재로그인·재부팅 시 설정이 초기화되는 사용자 환경에 부적합
- 기존 매핑 도구가 실행 중이면 중복 적용 방지를 위해 먼저 해제
- 이전 Kaymap 실행 중이면 종료 후 새 EXE 실행, 새 창 제목의 `CapsLock ↔ Ctrl` 표시 확인
- Kaymap 적용 중 Ctrl 조합은 물리 CapsLock으로 입력, 예: Parsec 토글 `CapsLock+Shift+I`

## 키 매핑

| Windows 입력 | Windows 출력 | Parsec을 거친 Mac 입력 |
| --- | --- | --- |
| Left Alt | Left Win | Left Command |
| Right Alt | Right Win | Right Command |
| Left Win | Left Alt | Left Option |
| Right Win | Right Alt | Right Option |
| CapsLock | Left Ctrl | Left Control |
| Left Ctrl | CapsLock | CapsLock |
| Right Ctrl | CapsLock | CapsLock |
| `\` | Backspace | 커서 왼쪽 문자 삭제 |
| Backspace | `\` 키 | `\` |
| Shift+Backspace | Shift+`\` 키 | US 배열 기준 `|` |

- Shift 및 매핑 대상 외 나머지 키: 원본 이벤트 유지
- Ctrl 반복 입력 및 좌우 Ctrl 동시 입력: 하나의 CapsLock 누름으로 처리, 마지막 Ctrl을 뗄 때 CapsLock 해제
- CapsLock 토글 상태: 시작·종료 시 별도 초기화 없음
- 반복 삭제: `\`를 누르고 있는 동안 반복 Backspace 전송
- 매핑 기준: Windows가 인식하는 좌우 Alt·Win·Ctrl, CapsLock, Backspace, `VK_OEM_5`
- 출력 기준: US ANSI 스캔 코드, HHKB 일본어 배열이나 다른 키보드 배열은 추가 확인 필요
- 한글 Windows 설정에서 오른쪽 Alt가 한영 키로 인식되는 경우 해당 키는 Right Alt 매핑 대상에 포함되지 않을 수 있음

## 실제 Windows 확인 항목

- [ ] 실행 후 작은 상태 창과 일시정지·종료 버튼 표시
- [ ] 로컬 메모장에서 `\`로 삭제, Backspace로 `\` 입력
- [ ] Mac에서 Alt+C·Alt+V·Alt+Z로 복사·붙여넣기·실행 취소
- [ ] Mac Terminal에서 물리 CapsLock+C로 Control+C 전달 및 실행 중 작업 중단
- [ ] 좌우 Ctrl을 각각 눌러 CapsLock 전환, 일반 영문 입력의 대소문자 전환 확인
- [ ] Ctrl을 길게 누르거나 좌우 Ctrl을 겹쳐 눌러 CapsLock 중복 전환 여부 확인
- [ ] Win+방향키로 Option+방향키 전달
- [ ] `\`를 길게 눌러 반복 삭제
- [ ] Shift+Backspace로 `|` 입력
- [ ] Alt+Tab으로 원격 Mac 앱 전환
- [ ] 일시정지 후 Windows 원래 키 배열 복구
- [ ] 종료 후 Windows 원래 키 배열 복구

## 호환성 범위

- 구현: Windows `WH_KEYBOARD_LL`로 원본 이벤트 차단, `SendInput`으로 교환한 스캔 코드 전송
- 자체 생성 입력의 재변환 방지, 종료·일시정지 시 눌린 출력 키 해제
- 키를 누른 상태에서 실행한 경우 일시정지로 시작, 모든 대상 키를 뗀 뒤 다시 적용
- `SendInput` 실패 시 자동 일시정지 및 상태 메시지 표시
- 프로그램의 전송 성공: Windows 입력 큐 삽입 성공 의미, Parsec이 원격 Mac에 전달했음을 보장하지 않음
- 로컬 메모장에서는 교환되지만 Parsec에서는 원래 키가 전달되는 경우: 이 방식과 해당 Parsec 환경의 호환성 미확보, 자동 전달 성공으로 간주 불가
- Parsec 실행 후 Kaymap 실행 순서로 시험, Parsec 재시작 후에는 Kaymap도 종료·재실행하여 시험
- 관리자 권한으로 실행된 대상에는 같은 권한 수준이 필요할 수 있음, 자동 권한 상승 없음
- 강제 프로세스 종료·OS 보안 화면·예약 시스템 단축키는 정상 종료 처리와 별개
- 직접 빌드한 서명되지 않은 EXE, 피씨방 실행 정책에 따른 차단 가능

## 개발 및 검증

- 언어: Go, 외부 모듈 의존성 없음
- TDD 기록: 동작 테스트 작성 → 미구현 식별자 오류로 실패 확인 → 구현 후 통과
- 자동 테스트 범위: 9개 키 매핑, Shift 보존, CapsLock의 Control 조합, 좌우 Ctrl의 CapsLock 공유, 반복 입력, 주입 입력 재귀 방지, 눌린 키 정리, 전송 실패·재시도, Windows INPUT 직렬화
- 검증 결과: `go test -race -cover ./...` 통과, remap 98.6%·wininput 100% 문장 커버리지, Windows 대상 `go vet` 통과
- Windows 어댑터: macOS에서 Windows x64 대상으로 교차 컴파일
- 실제 Windows UI·Hook 실행·Parsec 전달 검증: 미수행, 위 수동 항목으로 확인 필요

```sh
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags='-H=windowsgui -s -w' \
  -o dist/kaymap-windows-x64.exe \
  ./cmd/kaymap
```

## 근거 및 조사 기록

- Parsec 기본 Modifier 전달 및 Command·Ctrl 교환 — [Parsec 공식 문서](https://support.parsec.app/hc/en-us/articles/32361367389972-Swap-Command-and-Ctrl-for-MacOS)
- 원격 특수 단축키 전달 — [Parsec Immersive Mode](https://support.parsec.app/hc/en-us/articles/32361385571860-Immersive-Mode-Setting)
- Hook 처리 규칙 — [Microsoft LowLevelKeyboardProc](https://learn.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc)
- 합성 입력 및 권한 제약 — [Microsoft SendInput](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput)
- 기존 프로그램 비교 — [초기 조사 기록](docs/research.md), 피씨방 초기화 조건 반영 전 문서

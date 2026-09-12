# Kaymap — Windows 포터블 키 매핑

- 설치·재부팅 없이 실행하는 Windows 10/11 x64용 키 매핑 프로그램
- [EXE 다운로드](dist/kaymap-windows-x64.exe)
- [배포 ZIP 다운로드](dist/kaymap-windows-x64.zip), 실행 파일·사용법·Go 라이선스 포함
- 적용 범위: 실행 중 Windows 전체 키보드
- 일시정지·다시 적용·종료 및 입력 재연결 지원
- Registry 변경·드라이버 설치·자동 시작 등록 없음
- 별도 런타임 설치 불필요
- 현재 버전: 0.3, HHKB 모드 체크박스·입력 재연결·진단 표시 지원
- 검증 상태: 자동 테스트 및 Windows 빌드 완료, 실제 Windows UI·입력 전달은 추가 확인 필요

## 사용 절차

1. Windows PC로 EXE 복사
2. 기존 Kaymap이 실행 중이면 종료
3. 모든 키를 뗀 상태에서 `kaymap-windows-x64.exe` 실행
4. `HHKB 모드 — CapsLock ↔ Ctrl` 체크 여부 선택, 기본값 활성화
5. `적용 중` 표시 확인 후 사용할 창 클릭
6. 사용 종료 시 `종료` 버튼 클릭 또는 Kaymap 창 닫기

- EXE와 함께 설치할 추가 파일 없음
- 파일이 초기화되는 환경에서는 다음 사용 시 EXE만 다시 복사하여 실행
- 다른 키 매핑 도구가 실행 중이면 중복 적용 방지를 위해 먼저 해제
- HHKB 모드 활성화 시 Ctrl 조합은 물리 CapsLock으로 입력, 예: `CapsLock+C → Ctrl+C`

## HHKB 모드

- HHKB 사용자를 위한 CapsLock·Ctrl 교환 옵션, 실행 시 기본 활성화
- 체크 해제 시 Ctrl·CapsLock 원본 입력 유지, Alt·Win 및 역슬래시·Backspace 교환은 계속 적용

| Windows 입력 | 활성화 | 비활성화 |
| --- | --- | --- |
| CapsLock | Left Ctrl | CapsLock |
| Left Ctrl | CapsLock | Left Ctrl |
| Right Ctrl | CapsLock | Right Ctrl |

- 비활성화 시 Windows Ctrl을 변환하지 않으므로, 원격 프로그램이 Control을 그대로 전달하는 설정에서는 macOS Control로 입력
- 매핑 대상 키를 모두 뗀 뒤 옵션 변경, 누르는 중에는 변경을 거부하고 기존 체크 상태 유지
- 일시정지 상태에서 옵션 변경 가능, 전체 매핑은 일시정지 상태 유지
- 옵션은 현재 실행 중에만 유지, 프로그램 재실행 시 기본 활성화
- 비활성화 상태의 Ctrl 진단: `Ctrl(L) → Ctrl(L) / 원본 입력 유지`

## 키 매핑

- 아래 표는 HHKB 모드 활성화 기준

| Windows 입력 | Windows 출력 |
| --- | --- |
| Left Alt | Left Win |
| Right Alt | Right Win |
| Left Win | Left Alt |
| Right Win | Right Alt |
| CapsLock | Left Ctrl |
| Left Ctrl | CapsLock |
| Right Ctrl | CapsLock |
| `\` | Backspace |
| Backspace | `\` 키 |
| Shift+Backspace | US 배열 기준 `|` |

- Shift 및 매핑 대상 외 나머지 키: 원본 이벤트 유지
- HHKB 모드 활성화 시 좌우 Ctrl 동시 입력: 하나의 CapsLock 누름으로 처리, 마지막 Ctrl을 뗄 때 출력 키 해제
- CapsLock 토글 상태: 시작·종료 시 별도 초기화 없음
- 반복 삭제: `\`를 누르고 있는 동안 반복 Backspace 전송
- 매핑 기준: Windows가 인식하는 좌우 Alt·Win·Ctrl, CapsLock, Backspace, `VK_OEM_5`
- 출력 기준: US ANSI 스캔 코드, 다른 키보드 배열은 추가 확인 필요
- 한글 Windows 설정에서 오른쪽 Alt가 한영 키로 인식되는 경우 Right Alt 매핑 대상에 포함되지 않을 수 있음

## 입력 진단 및 재연결

- 창 전환 후 0.5초 이상 안정된 상태에서 입력 연결 갱신
- 대상 키가 눌려 있거나 일시정지 상태면 갱신 연기
- `입력 다시 연결`: 적용 중 버튼 클릭 후 사용할 창으로 돌아가면 2초 이후 연결 갱신
- 마지막 매핑 키와 Windows 전송 결과 표시, 일반 문자 입력은 기록하지 않음
- 진단 정보는 메모리에만 유지, 파일 저장 및 네트워크 전송 없음

| 진단 표시 | 의미 |
| --- | --- |
| Win(L) → Alt(L) / Windows 전송 성공 | Win 감지 및 Alt 입력 큐 삽입 성공, 대상 앱의 수신 확인과 별개 |
| 입력 대기 또는 입력 번호가 갱신되지 않음 | 해당 입력이 Kaymap에 도달하지 않음, 적용 상태 및 입력 연결 확인 필요 |
| Windows 전송 실패 | Windows 입력 삽입 실패, 대상 앱과의 권한 수준 확인 필요 |
| Alt(L) → Win(L) | Windows가 해당 물리 키를 Alt로 인식, 키보드 모드·다른 매핑 확인 필요 |

## 확인 항목 및 제약

- [ ] 메모장에서 `\`로 삭제, Backspace로 `\` 입력
- [ ] `CapsLock+A`로 Ctrl+A 전달 및 전체 선택
- [ ] 좌우 Ctrl을 각각 눌러 CapsLock 전환
- [ ] Ctrl을 길게 누르거나 겹쳐 눌러 CapsLock 중복 전환 여부 확인
- [ ] HHKB 모드 해제 후 Ctrl+A로 전체 선택, CapsLock으로 대소문자 전환
- [ ] HHKB 모드를 다시 켰을 때 CapsLock·Ctrl 교환 복구
- [ ] 키를 누른 상태의 옵션 변경 거부 및 해제 후 변경 가능 여부 확인
- [ ] 진단 표시에서 `Win → Alt`, `Alt → Win` 확인
- [ ] Shift+Backspace로 `|` 입력
- [ ] 일시정지·종료 후 원래 키 배열 복구

- 원격 접속 시 최종 입력은 원격 프로그램과 대상 OS의 변환 설정에 따라 달라짐
- 일부 원격 환경에서 Win 입력이 원하는 결과로 변환되지 않는 사용자 보고 존재, 입력 재연결 이후 개선 여부는 실기 확인 필요
- 키보드 Hook을 거치지 않는 입력을 사용하는 앱의 호환성은 보장하지 않음
- 관리자 권한으로 실행된 대상에는 같은 권한 수준이 필요할 수 있음, 자동 권한 상승 없음
- 강제 종료·OS 보안 화면·예약 시스템 단축키는 정상 종료 처리와 별개
- 서명되지 않은 EXE, PC의 실행 정책에 따른 차단 가능

## 개발 및 검증

- 언어: Go, 외부 모듈 의존성 없음
- 구현: `WH_KEYBOARD_LL`로 원본 이벤트 차단, `SendInput`으로 교환한 스캔 코드 전송
- 자체 생성 입력의 재변환 방지, 종료·일시정지 시 눌린 출력 키 해제
- TDD: 테스트 작성 → 실패 확인 → 구현 후 통과
- 자동 테스트: 9개 키 매핑·HHKB 기본값과 옵션 전환·원본 Ctrl 유지·반복 입력·Modifier 조합·입력 정리·전송 실패·INPUT 직렬화·재연결 시점·진단 표시
- 실제 Windows UI·Hook 동작·대상 앱 전달은 추가 검증 필요

```sh
go test -race -cover ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go vet ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags='-H=windowsgui -s -w' \
  -o dist/kaymap-windows-x64.exe \
  ./cmd/kaymap
```

## 참고 자료

- [Microsoft LowLevelKeyboardProc](https://learn.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc)
- [Microsoft SendInput](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput)
- [Microsoft Hooks Overview](https://learn.microsoft.com/en-us/windows/win32/winmsg/about-hooks)
- [기존 프로그램 비교](docs/research.md)

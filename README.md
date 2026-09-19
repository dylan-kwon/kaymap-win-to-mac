# Kaymap — Parsec용 Windows 키 매핑

- 설치·재부팅 없이 실행하는 Windows 10/11 x64용 포터블 프로그램
- 현재 버전: 0.10, 현재 키보드 실측에 맞춘 Command·Control 출력 교환
- [EXE 다운로드](dist/kaymap-windows-x64.exe)
- [배포 ZIP 다운로드](dist/kaymap-windows-x64.zip), 실행 파일·사용법·Go 라이선스 포함
- 실행 중 Windows 전체 키보드에 적용, 일시정지·종료 시 원본 입력 복구
- Registry 변경·드라이버 설치·자동 시작 등록·별도 런타임 설치 없음

## 전제 및 사용 절차

- 현재 사용 중인 Parsec 설정 유지, 키보드 실측 결과에 맞춰 출력 교환
- [Parsec Command·Control 교환 설정](https://support.parsec.app/hc/en-us/articles/32361367389972-Swap-Command-and-Ctrl-for-MacOS)
- 사용자가 보고한 CapsLock → Command·Alt → Control 동작에서 역산한 출력 변환 기준

| Windows에서 Parsec에 전달한 키 | Mac에서 감지되는 키 |
| --- | --- |
| Ctrl | Control |
| Alt | Option |
| Win | Command |

1. 기존 Kaymap 종료
2. 모든 키를 뗀 상태에서 새 `kaymap-windows-x64.exe` 실행
3. `Kaymap 0.10 — Parsec 매핑 ON` 창 확인
4. 필요한 경우 `HHKB 모드` 체크 변경, 기본값 활성화
5. `적용 중` 표시 확인 후 Parsec 화면 클릭
6. 물리 Alt+A 전체 선택·Alt+W 창 닫기·Win 키의 Option 인식 확인

- Alt·Win 전용 체크박스 없음, 현재 키보드 실측 기준 매핑 고정
- Parsec의 Command·Control 교환 설정 변경 시 최종 Mac 입력이 달라질 수 있으므로 현재 설정 유지 필요
- EXE만으로 실행 가능, 파일이 초기화된 환경에서는 다음 사용 시 EXE 복사 후 실행
- 다른 키 매핑 도구 사용 시 중복 변환 확인 필요

## Alt·Win 고정 매핑

| 물리 Windows 키 | Kaymap의 Windows 출력 | Parsec을 거친 Mac 입력 |
| --- | --- | --- |
| Left Alt | Left Win | Left Command |
| Right Alt | Right Win | Right Command |
| Left Win | Left Alt | Left Option |
| Right Win | Right Alt | Right Option |

- HHKB 모드와 무관하게 위 매핑 적용
- 재실행·입력 재연결·다시 적용 후에도 같은 매핑 사용
- 일시정지·종료 시 원본 입력 복구, Mac 입력은 위의 Parsec 기본 변환표로 복귀
- Alt와 Win 동시 입력 시 Command와 Option을 각각 유지
- 각 키의 누름과 해제를 쌍으로 처리, 자체 생성 입력의 재변환 방지

## HHKB 모드

- CapsLock·Ctrl의 Mac 역할 교환 및 Backspace·역슬래시 교환 옵션
- 기본 활성화, 변경 시 매핑 대상 키를 모두 뗀 상태 필요
- 현재 실행 중에만 설정 유지, 재실행 시 기본 활성화

| 물리 Windows 키 | HHKB On: Windows 출력 → Mac 입력 | HHKB Off: Windows 출력 → Mac 입력 |
| --- | --- | --- |
| CapsLock | Left Ctrl → Control | CapsLock → CapsLock |
| Left Ctrl | CapsLock → CapsLock | Left Ctrl → Control |
| Right Ctrl | CapsLock → CapsLock | Right Ctrl → Control |
| Backspace | `\` → `\` | Backspace → Backspace |
| `\` | Backspace → Backspace | `\` → `\` |

- HHKB Off에서도 Ctrl을 Windows Ctrl로 출력하여 Mac Control 유지
- Mac Command+C 복사: HHKB 설정과 무관하게 물리 Alt+C
- Mac Control+C: HHKB On이면 물리 CapsLock+C, Off이면 물리 Ctrl+C
- Windows 로컬 Ctrl+C 복사: HHKB On이면 물리 CapsLock+C, Off이면 물리 Ctrl+C
- HHKB On에서 좌우 Ctrl 동시 입력: CapsLock 한 번 누름으로 합치고 마지막 Ctrl 해제 시 출력 해제
- HHKB On에서 역슬래시를 누르고 있으면 Backspace 반복 전송
- Shift+Backspace: HHKB On 및 US ANSI 배열 기준 `|` 입력
- Shift 및 매핑 대상 외 키: 원본 입력 유지
- CapsLock 토글 상태: 시작·종료 시 별도 초기화 없음
- 일시정지 중 HHKB 변경 가능, 옵션 변경만으로 매핑 재개 없음

## 입력 진단 및 재연결

- 창 전환 후 0.5초 이상 안정된 상태에서 입력 연결 갱신
- 대상 키를 누르고 있거나 일시정지 중이면 갱신 연기
- `입력 다시 연결`: 적용 중 버튼 클릭 후 대상 창으로 돌아가면 2초 이후 연결 갱신
- 마지막 매핑 대상 키와 Windows 전송 결과만 메모리에 표시
- 일반 문자 입력 기록·파일 저장·네트워크 전송 없음

| 진단 표시 예시 | 의미 |
| --- | --- |
| Alt(L) → Win(L) / Windows 전송 성공 | Command용 Win 입력 큐 삽입 성공 |
| Win(L) → Alt(L) / Windows 전송 성공 | Option용 Alt 입력 큐 삽입 성공 |
| CapsLock → Ctrl(L) / Windows 전송 성공 | HHKB On에서 Control용 Ctrl 입력 큐 삽입 성공 |
| Ctrl(L) → Ctrl(L) / Windows 전송 성공 | HHKB Off에서 Control용 Ctrl 입력 큐 삽입 성공 |
| Windows 전송 실패 | Windows 입력 삽입 실패, 대상 앱과의 권한 수준 확인 필요 |
| 입력 대기 또는 입력 번호 미갱신 | 해당 입력의 Hook 도달 여부 및 적용 상태 확인 필요 |

- Windows 전송 성공은 Parsec이나 Mac 앱의 실제 수신 확인과 별개

## 검증 및 제약

- TDD: 변경된 키 매핑 테스트 실패 확인 후 구현 및 통과
- 자동 검증: 보고된 Parsec 변환표와 Kaymap 출력을 합성한 최종 Mac 역할 확인
- 자동 검증: 좌우 Modifier·A/W 조합·동시 입력 및 해제 순서·일시정지/재개·HHKB 전환·반복 입력·입력 정리·전송 실패·INPUT 직렬화·입력 재연결·진단 표시
- Windows x64 빌드·배포 ZIP 무결성 검증 완료
- 실제 Windows Hook·Parsec·Mac 앱을 통한 0.10의 Alt+W 동작은 실기 확인 필요
- 출력 기준: US ANSI 스캔 코드, 다른 키보드 배열 및 한영 전용 키는 추가 확인 필요
- Hook을 거치지 않는 입력을 사용하는 앱의 호환성은 보장하지 않음
- 관리자 권한 대상 앱에는 같은 권한 수준이 필요할 수 있음, 자동 권한 상승 없음
- 강제 종료·OS 보안 화면·예약 시스템 단축키는 정상 종료 처리와 별개
- 서명되지 않은 EXE, PC의 실행 정책에 따른 차단 가능

## 개발

- Go 사용, 외부 모듈 의존성 없음
- `WH_KEYBOARD_LL`로 원본 이벤트 처리, `SendInput`으로 변환한 Windows 스캔 코드 전송

```sh
go test -race -cover ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go vet ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags='-H=windowsgui -s -w' \
  -o dist/kaymap-windows-x64.exe \
  ./cmd/kaymap
```

- [Microsoft LowLevelKeyboardProc](https://learn.microsoft.com/en-us/windows/win32/winmsg/lowlevelkeyboardproc)
- [Microsoft SendInput](https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendinput)
- [기존 프로그램 비교](docs/research.md)

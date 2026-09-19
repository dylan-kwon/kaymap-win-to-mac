# Kaymap — Parsec용 Windows 키 매핑

- 설치·재부팅 없이 실행하는 Windows 10/11 x64용 포터블 프로그램
- 현재 버전: 1.0.0
- [EXE 다운로드](https://github.com/dylan-kwon/kaymap-win-to-mac/releases/download/1.0.0/kaymap-windows-x64.exe)
- [배포 ZIP 다운로드](https://github.com/dylan-kwon/kaymap-win-to-mac/releases/download/1.0.0/kaymap-windows-x64.zip), 실행 파일·사용법·Go 라이선스 포함
- 실행 중 Windows 전체 키보드에 적용, 일시정지·종료 시 원본 입력 복구
- Registry 변경·드라이버 설치·자동 시작 등록·별도 런타임 설치 없음

## 전제 및 사용 절차

- 목표 Mac 입력: Alt → Command, CapsLock → Control(HHKB On)
- 아래 표는 Windows Ctrl → Mac Control, Windows Win → Mac Command 변환을 전제로 한 예상 입력
- 기존 Parsec 설정 유지, 설정이나 키보드 환경이 다른 경우 실제 Mac 입력 확인 필요
- [Parsec Command·Control 교환 설정](https://support.parsec.app/hc/en-us/articles/32361367389972-Swap-Command-and-Ctrl-for-MacOS)

| Windows에서 Parsec에 전달한 키 | 예상 Mac 입력 |
| --- | --- |
| Ctrl | Control |
| Alt | Option |
| Win | Command |

1. 기존 Kaymap 종료
2. 모든 키를 뗀 상태에서 새 `kaymap-windows-x64.exe` 실행
3. `Kaymap 1.0.0 — Parsec 매핑 ON` 창 확인
4. 필요한 경우 `HHKB 모드` 체크 변경, 기본값 활성화
5. `적용 중` 표시 확인 후 Parsec 화면 클릭
6. 물리 Alt → Command·CapsLock → Control·Win → Option 인식 확인

- Alt·Win 전용 체크박스 없음, 현재 키보드 실측 기준 매핑 고정
- Parsec의 Command·Control 교환 설정 변경 시 최종 Mac 입력이 달라질 수 있으므로 현재 설정 유지 필요
- EXE만으로 실행 가능, 파일이 초기화된 환경에서는 다음 사용 시 EXE 복사 후 실행
- 다른 키 매핑 도구 사용 시 중복 변환 확인 필요

## Alt·Win 고정 매핑

| 물리 Windows 키 | Kaymap의 Windows 출력 | 예상 Mac 입력 |
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

## 문제 해결 및 개발 문서

- [입력 진단·재연결·제약 사항](docs/troubleshooting.md)
- [개발·검증·배포 파일](docs/development.md)
- [1.0.0 릴리즈 노트](docs/releases/1.0.0.md)
- [초기 도구 검토 기록](docs/research.md)

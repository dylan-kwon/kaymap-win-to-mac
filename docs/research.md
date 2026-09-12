# Windows → Mac 원격 접속 키 매핑 검토

- 보관 문서: 피씨방 재부팅 시 설정 초기화 조건을 반영하기 전 조사, 현재 사용법은 [README](../README.md) 참조
- 검토일: 2026-09-12
- 결론: Parsec 기본 전달 설정과 Windows 키 교환 조합으로 구현 가능, 도구별 실제 원격 입력 검증 필요
- 사용 환경: Windows 클라이언트 → Parsec → Mac 호스트
- 현재 상태: Parsec 공식 매핑 확인 완료, 실행 코드 및 Windows 실기 검증 미수행

## Parsec 설정 및 권장 순서

1. Windows의 Parsec → Settings → Client → `Swap Command and Ctrl for MacOS`를 `Off`로 설정
2. 원격 세션에서 Keyboard Immersive Mode 활성화, Windows 기본 토글 단축키 `Ctrl+Shift+I` 사용 가능
3. 아래 표대로 Windows의 Alt·Win 및 `\`·Backspace 교환 적용
4. Mac에서 복사·붙여넣기, Terminal의 Ctrl+C, 반복 삭제, Shift+Backspace, 앱 전환 검증

- Parsec 기본 전달: Windows `Win → Command`, `Alt → Option`, `Ctrl → Control` — [Parsec 공식 Modifier 매핑 문서](https://support.parsec.app/hc/en-us/articles/32361367389972-Swap-Command-and-Ctrl-for-MacOS)
- Keyboard Immersive Mode: Alt+Tab 등 특수 키 조합을 원격 호스트로 전달하는 기능, 키 교환 기능과 별개 — [공식 Immersive Mode 문서](https://support.parsec.app/hc/en-us/articles/32361385571860-Immersive-Mode-Setting), [공식 단축키 문서](https://support.parsec.app/hc/en-us/articles/32381778420372-Configure-Hotkeys)
- 공식 문서에서 확인한 내장 키 교환: Command·Ctrl 교환, 요구한 Alt·Win 및 Backspace 교환 전체를 위한 내장 기능은 확인되지 않음
- Windows 전체 키 배열 변경 허용 시: SharpKeys를 통한 OS 키 교환 검토, Parsec 입력 결과는 실기 확인 필요
- 접속 전후 즉시 전환 필요 시: PowerToys 또는 AutoHotkey의 소규모 호환성 시험 우선
- Parsec 창에서 PowerToys 단축키 매핑이 작동하지 않는 과거 보고 존재 — [PowerToys 공식 저장소 이슈 #21515](https://github.com/microsoft/PowerToys/issues/21515)
- 해당 보고: 2022년 PowerToys 0.63.0 단축키 매핑 사례, 현재 버전의 모든 단일 키 매핑 실패 또는 AutoHotkey 실패를 입증하는 자료는 아님
- 별도 개발 여부: 즉시 전환 필요성과 실제 입력 시험 결과를 기준으로 결정, 일반 Hook 구현만으로 호환성이 해결된다고 가정하지 않음

## 요구사항

| Windows에서 누르는 키 | Mac에서 원하는 입력 |
| --- | --- |
| Alt | Command |
| Win | Option |
| Ctrl | Control |
| `\` | Backspace — 커서 왼쪽 문자 삭제 |
| Backspace | `\` |

- Mac의 Backspace: 일반적인 Mac 키보드의 `delete` 기능, 앞쪽 문자를 지우는 Forward Delete와 구분
- 키 이름 기준: 키캡의 인쇄보다 Windows에서 실제 인식하는 입력 기준
- HHKB 모델·동작 모드·DIP 설정에 따른 실제 입력 확인 필요

## 변환 구조

- 입력 경로: 물리 키 → Windows 키 매핑 → 원격 프로그램의 키 변환 → Mac 입력
- Parsec의 Command·Ctrl 교환을 끈 상태에서 사용할 Windows 설정:

| Windows 입력 | Windows 매핑 출력 | 예상 Mac 입력 |
| --- | --- | --- |
| Left Alt | Left Win | Left Command |
| Right Alt | Right Win | Right Command |
| Left Win | Left Alt | Left Option |
| Right Win | Right Alt | Right Option |
| Ctrl | 변경 없음 | Control |
| `\` | Backspace | Backspace |
| Backspace | `\` 키 | `\` |

- 위 표: Parsec 공식 기본 매핑에 근거한 설정안, Windows 매핑 도구와 Parsec 간 호환성은 실기 검증 필요
- Parsec의 `Swap Command and Ctrl for MacOS` 사용 시 `Ctrl → Control` 요구사항과 충돌하므로 비활성화
- `Backspace → \`: 문자열 붙여넣기보다 키 자체 교환 권장, Shift 조합과 길게 누르기 동작 검증 필요

## 기존 프로그램 비교

| 방법 | 적합한 상황 | 주요 제약 | 판단 |
| --- | --- | --- | --- |
| Parsec 내장 키 변환 | 기본 Modifier 전달 및 특수 조합 전달 설정 | 요구한 전체 키 교환 기능 미확인 | Command·Ctrl 교환 Off, Keyboard Immersive Mode 활성화 |
| PowerToys Keyboard Manager | GUI로 간단히 설정하고 원격 작업 전후 켜고 끄는 경우 | 전체 적용, 상주 필요, Parsec 관련 과거 문제 보고 | 소규모 호환성 시험용 |
| AutoHotkey v2 | 원격 창 활성 상태에 따라 자동 적용하는 경우 | 스크립트 작성 및 Parsec 입력 호환성 검증 필요 | 자동 전환 필요 시 시험 후보 |
| SharpKeys | Windows 전체 키 배열 변경을 허용하는 경우 | Registry 기반, 로그아웃 또는 재부팅 필요, 앱별 전환 불가 | 전체 적용 허용 시 검토 후보 |

- PowerToys: 키 교환과 앱별 단축키 지원, 앱별 단축키 설정과 단일 키 전체 교환의 적용 범위 구분 필요 — [Microsoft 공식 문서](https://learn.microsoft.com/en-us/windows/powertoys/keyboard-manager)
- AutoHotkey: 조건별 매핑, Modifier 유지, 반복 입력 지원 — [AutoHotkey 공식 Remap 문서](https://github.com/AutoHotkey/AutoHotkeyDocs/blob/v2/docs/misc/Remap.htm)
- AutoHotkey 배포: Ahk2Exe로 스크립트를 Windows 실행 파일로 패키징 가능 — [공식 Ahk2Exe 저장소](https://github.com/AutoHotkey/Ahk2Exe)
- SharpKeys: Windows Registry의 키 매핑 관리, 상주 프로그램 없이 적용 — [공식 SharpKeys 저장소](https://github.com/randyrants/sharpkeys)

## SharpKeys 설정 절차

1. [공식 SharpKeys 배포 페이지](https://github.com/randyrants/sharpkeys/releases)에서 설치 파일 다운로드 및 Windows에 설치
2. 기존 매핑이 있으면 현재 목록을 기록하고 Save Keys로 보관
3. Add에서 아래 Windows 매핑 6개 등록, Alt는 Type Key 대신 목록에서 선택
4. `Left Alt → Left Windows`, `Left Windows → Left Alt` 등록
5. `Right Alt → Right Windows`, `Right Windows → Right Alt` 등록
6. `\ → Backspace`, `Backspace → \` 등록
7. Ctrl 매핑은 추가하지 않음
8. Write to Registry 선택 후 로그아웃·재로그인 또는 재부팅
9. Parsec 설정 확인 후 원격 Mac에서 검증

- 적용 범위: Windows 전체 및 모든 사용자, HHKB만 선택 적용하거나 Parsec 종료 시 자동 해제 불가
- 복구: 이번에 추가한 행만 제거하거나 기존 목록 복원 → Write to Registry → 로그아웃·재로그인 또는 재부팅
- 근거: [SharpKeys 공식 사용법 및 제약](https://github.com/randyrants/sharpkeys)

## PowerToys 시험 절차

1. Windows에 [PowerToys 설치](https://learn.microsoft.com/en-us/windows/powertoys/install)
2. Keyboard Manager → Remap a key 진입
3. 위 Windows 매핑 표의 Alt·Win 교환 및 `\`·Backspace 교환 등록
4. Ctrl은 별도 매핑 없이 유지
5. 원격 프로그램의 자동 단축키 변환 및 원격 단축키 전달 설정 확인
6. 아래 검증 항목 수행
7. 원격 작업 종료 후 Keyboard Manager 비활성화 또는 PowerToys 종료

- 여러 매핑 도구를 동시에 적용할 경우 중복 변환 가능, 한 가지 방식씩 검증
- PowerToys 공식 제약: `Win+L`과 `Ctrl+Alt+Del` 등 OS 예약 조합 변경 불가, 관리자 권한 창과 키보드 Hook 사용 앱에서 동작 제약 가능 — [Microsoft 공식 문서](https://learn.microsoft.com/en-us/windows/powertoys/keyboard-manager)

## 별도 개발 시 범위

- 즉시 전환이 필요한 경우의 시험 구현 후보: AutoHotkey v2 기반 Windows 트레이 프로그램
- 대상 원격 창이 활성화된 동안 매핑 적용
- 트레이 메뉴의 적용·해제 및 종료 기능
- 좌우 Alt·Win 개별 교환, Ctrl 유지, `\`·Backspace 교환
- 키를 누르는 시점과 떼는 시점 추적, 반복 입력 및 Modifier 조합 보존
- 키를 누른 채 창 전환·매핑 해제·종료하는 경우 눌림 상태 정리
- 브라우저 기반 원격 접속 사용 시 대상 세션 구분 필요, 브라우저 프로세스 전체에 적용 시 다른 탭에도 영향 가능
- 개발 순서: 동작 테스트 작성 → 실패 확인 → 구현 → 테스트 통과 → 실제 Windows–Mac 원격 접속 검증
- 네이티브 Hook 기반 프로그램 검토 조건: 기존 도구로 해결되지 않는 구체적인 입력 호환성 문제 확인 시

## 실제 환경 검증 항목

| Windows에서 수행할 동작 | Mac에서 기대하는 결과 |
| --- | --- |
| Alt+C / Alt+V / Alt+Z | Command+C / Command+V / Command+Z |
| Alt+Shift+Z | Command+Shift+Z |
| Win+좌우 방향키 | Option+좌우 방향키 |
| Terminal에서 Ctrl+C | Control+C로 실행 중 작업 중단 |
| `\` 입력 및 길게 누르기 | 커서 왼쪽 문자 삭제 및 반복 삭제 |
| Backspace | `\` 입력 |
| Shift+Backspace | US 배열 기준 `|` 입력 |
| Alt+Tab | 의도한 원격 Mac 앱 전환, 로컬 Windows 가로채기 여부 확인 |
| Alt를 누른 채 창 전환 후 떼기 | Command 눌림 상태가 남지 않음 |
| 원격 창 종료 또는 매핑 해제 | 로컬 Windows에서 원래 입력 복구 |
| 한글·영문 입력 상태에서 동일 시험 | 지정 키 유지 및 한영 전환 충돌 없음 |

## 설정 확정에 필요한 정보

- Parsec 버전
- HHKB 모델과 현재 모드, Windows에서 인식하는 키
- 원격 창에서만 적용할지, 수동으로 전체 적용을 전환할지
- 현재 검증 한계: macOS 작업 환경에서 문서 조사만 수행, Windows 실행 및 원격 Mac 입력 결과 미검증

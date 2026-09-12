# Windows 키 매핑 도구 검토

- 검토일: 2026-09-12
- 목적: 설치·재부팅 없이 사용할 수 있는 키 매핑 방식 검토
- 현재 구현 및 사용법: [README](../README.md) 참조

## 기존 프로그램 비교

| 방법 | 적합한 상황 | 주요 제약 |
| --- | --- | --- |
| PowerToys Keyboard Manager | GUI로 키 교환 설정 및 수동 적용 전환 | 단일 키 매핑은 전체 적용, 백그라운드 실행 필요 |
| AutoHotkey v2 | 창 활성 상태에 따른 조건부 매핑 | 스크립트 작성 및 대상 앱과의 입력 호환성 검증 필요 |
| SharpKeys | Windows 전체 키 배열을 지속적으로 변경 | Registry 기반, 로그아웃 또는 재부팅 필요, 앱별 전환 불가 |
| Windows Hook 기반 프로그램 | 설치·재부팅 없이 키 교환 및 진단 제공 | 대상 앱의 입력 처리 방식·권한 수준에 따른 호환성 확인 필요 |

- PowerToys: 키 교환 및 앱별 단축키 지원 — [Microsoft 공식 문서](https://learn.microsoft.com/en-us/windows/powertoys/keyboard-manager)
- AutoHotkey: 조건별 매핑, Modifier 유지, 반복 입력 지원 — [공식 Remap 문서](https://github.com/AutoHotkey/AutoHotkeyDocs/blob/v2/docs/misc/Remap.htm)
- AutoHotkey 배포: 스크립트를 Windows 실행 파일로 패키징 가능 — [공식 Ahk2Exe 저장소](https://github.com/AutoHotkey/Ahk2Exe)
- SharpKeys: Windows Registry 키 매핑 관리 — [공식 저장소](https://github.com/randyrants/sharpkeys)

## 구현 선택 및 검증 기준

- 선택: Windows Hook 기반 포터블 프로그램
- 적용·일시정지·종료를 실행 중 전환, 재부팅 불필요
- 키 누름·뗌·반복 입력·Modifier 조합 보존
- 일시정지·종료 시 눌린 출력 키 정리
- 재부팅 후 설정이 초기화되는 환경에서도 EXE 재실행으로 적용
- 다른 매핑 도구와 중복 적용 여부 확인
- 원격 환경에서는 로컬 키 교환과 원격 프로그램의 추가 변환을 구분하여 검증
- 검증 범위: 로컬 입력 큐 삽입 성공과 대상 앱의 실제 수신을 별도로 확인
- 일반 Hook 구현만으로 모든 앱과의 호환성이 해결된다고 가정하지 않음

# 개발 및 검증

- Go 1.24 이상 사용, 외부 모듈 의존성 없음
- 배포 대상: Windows x64
- 키 매핑 로직 변경 시 TDD 적용: 테스트 실패 확인 → 구현 → 통과 확인
- 버전 변경 시 창 제목·README·빠른 시작 안내·릴리즈 노트 동기화

## 빌드

- `WH_KEYBOARD_LL`로 원본 이벤트 처리, `SendInput`으로 Alt는 가상 키, 나머지는 Windows 스캔 코드 전송

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
- [기존 프로그램 비교](research.md)

## 배포 파일

- `dist/kaymap-windows-x64.exe`: Windows 실행 파일
- `dist/QUICK_START.txt`: 빠른 시작 안내
- `dist/THIRD_PARTY_NOTICES.txt`: Go 라이선스 고지
- `dist/SHA256SUMS.txt`: 실행 파일 및 ZIP의 SHA-256 체크섬
- `dist/kaymap-windows-x64.zip`: 실행 파일·빠른 시작 안내·라이선스 포함
- 빌드 후 ZIP 재생성, 압축 무결성 및 내부 실행 파일 일치 확인
- ZIP 생성 후 EXE·ZIP 체크섬 갱신, GitHub Release에 EXE·ZIP·체크섬 첨부

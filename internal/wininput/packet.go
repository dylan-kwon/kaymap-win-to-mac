// Package wininput은 Windows 64비트 INPUT 구조체를 직렬화한다.
package wininput

import "encoding/binary"

const Marker uint64 = 0x4B41594D4150

// Encode는 US ANSI 키 위치의 스캔 코드를 사용한다.
// Windows x64/ARM64의 INPUT 공용체 정렬과 크기를 명시적으로 보존한다.
func Encode(key uint32, down bool) ([40]byte, bool) {
	var packet [40]byte
	var scan uint16
	flags := uint32(0x0008)
	switch key {
	case 0xA4:
		scan = 0x38
	case 0xA5:
		scan = 0x38
		flags |= 0x0001
	case 0x5B:
		scan = 0x5B
		flags |= 0x0001
	case 0x5C:
		scan = 0x5C
		flags |= 0x0001
	case 0x08:
		scan = 0x0E
	case 0xDC:
		scan = 0x2B
	case 0xA2:
		scan = 0x1D
	case 0xA3:
		scan = 0x1D
		flags |= 0x0001
	case 0x14:
		scan = 0x3A
	default:
		return packet, false
	}
	if !down {
		flags |= 0x0002
	}
	binary.LittleEndian.PutUint32(packet[0:4], 1)
	binary.LittleEndian.PutUint16(packet[10:12], scan)
	binary.LittleEndian.PutUint32(packet[12:16], flags)
	binary.LittleEndian.PutUint64(packet[24:32], Marker)
	return packet, true
}

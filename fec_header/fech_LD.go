package fech

import (
	"encoding/binary"
)

type FecHeaderLD struct {
	R                 bool
	F                 bool
	P                 bool
	X                 bool
	CC                uint8
	M                 bool
	PTRecovery        uint8
	LengthRecovery    uint16
	TimestampRecovery uint32
	SN_base           uint16
	L                 uint8
	D                 uint8
}

/*

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |0|1|P|X|  CC   |M| PT recovery |         length recovery       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                          TS recovery                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           SN base_i           |  L (columns)  |    D (rows)   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |    ... next SN base and L/D for CSRC_i in CSRC list ...       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   :               Repair "Payload" follows FEC Header             :
   :                                                               :
*/

func NewFecHeaderLD(R bool, F bool, P bool, X bool, CC uint8, M bool, PTRecovery uint8, LengthRecovery uint16, TimestampRecovery uint32, SN_base uint16, L uint8, D uint8) FecHeader {
	return &FecHeaderLD{
		R:                 R,
		F:                 F,
		P:                 P,
		X:                 X,
		M:                 M,
		L:                 L,
		D:                 D,
		CC:                CC,
		SN_base:           SN_base,
		PTRecovery:        PTRecovery,
		LengthRecovery:    LengthRecovery,
		TimestampRecovery: TimestampRecovery,
	}
}

func (fh *FecHeaderLD) Marshal() []byte {
	// size to be detarmined later , for now for L and D variant size = 12 bytes
	size := 12

	buf := make([]byte, size)

	if fh.R {
		buf[0] = (1 << 7)
	}

	if fh.F {
		buf[0] |= (1 << 6)
	}

	if fh.P {
		buf[0] |= (1 << 5)
	}

	if fh.X {
		buf[0] |= (1 << 4)
	}

	buf[0] |= fh.CC

	if fh.M {
		buf[1] = (1 << 7)
	}

	buf[1] |= (fh.PTRecovery & 0x7F)

	binary.BigEndian.PutUint16(buf[2:4], fh.LengthRecovery)
	binary.BigEndian.PutUint32(buf[4:8], fh.TimestampRecovery)

	binary.BigEndian.PutUint16(buf[8:10], fh.SN_base)

	buf[10] = fh.L
	buf[11] = fh.D

	return buf
}

func (fh *FecHeaderLD) Unmarshal(buf []byte) {
	fh.R = (buf[0] >> 7 & 0x1) > 0
	fh.F = (buf[0] >> 6 & 0x1) > 0
	fh.P = (buf[0] >> 5 & 0x1) > 0
	fh.X = (buf[0] >> 4 & 0x1) > 0
	fh.CC = uint8((buf[0] & uint8(0xF)))
	fh.M = (buf[0] >> 7 & 0x1) > 0

	fh.PTRecovery = buf[1] & 0x7F

	fh.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])
	fh.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])

	fh.SN_base = binary.BigEndian.Uint16(buf[8:10])
	fh.L = buf[10]
	fh.D = buf[11]

}

package fech

import "encoding/binary"

type FecHeaderFlexibleMask struct {
	R                 bool
	F                 bool
	P                 bool
	X                 bool
	CC                uint8
	M                 bool
	PTRecovery        uint8 // use only 7 bits
	LengthRecovery    uint16
	TimestampRecovery uint32
	SN_base           uint16
	K1                bool
	Mask              uint16 // tail 15 bits only
	OptionalMask1     uint32 // tail 31 bits only
	K2                bool
	OptionalMask2     uint64 // use all bits
}

/*

 0                   1                   2                   3
      0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |0|0|P|X|  CC   |M| PT recovery |        length recovery        |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |                          TS recovery                          |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |           SN base_i           |k|          Mask [0-14]        |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |k|                   Mask [15-45] (optional)                   |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |                     Mask [46-109] (optional)                  |
     |                                                               |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |   ... next SN base and Mask for CSRC_i in CSRC list ...       |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     :               Repair "Payload" follows FEC Header             :
     :                                                               :
*/

func NewFecHeaderFlexibleMask(R bool, F bool, P bool, X bool, CC uint8, M bool, PTRecovery uint8, LengthRecovery uint16, TimestampRecovery uint32, SN_base uint16, K1 bool, Mask uint16, OptionalMask1 uint32, K2 bool, OptionalMask2 uint64) FecHeader {
	return &FecHeaderFlexibleMask{
		R:                 R,
		F:                 F,
		P:                 P,
		X:                 X,
		CC:                CC,
		M:                 M,
		PTRecovery:        PTRecovery,
		LengthRecovery:    LengthRecovery,
		TimestampRecovery: TimestampRecovery,
		SN_base:           SN_base,
		K1:                K1,
		Mask:              Mask,
		OptionalMask1:     OptionalMask1,
		K2:                K2,
		OptionalMask2:     OptionalMask2,
	}
}

func (ff *FecHeaderFlexibleMask) Marshal() []byte {

	size := 24 // to be determined later

	buf := make([]byte, size)

	if ff.R {
		buf[0] = (1 << 7)
	}

	if ff.F {
		buf[0] |= (1 << 6)
	}

	if ff.P {
		buf[0] |= (1 << 5)
	}

	if ff.X {
		buf[0] |= (1 << 4)
	}

	buf[0] |= ff.CC

	if ff.M {
		buf[1] = (1 << 7)
	}

	buf[1] |= (ff.PTRecovery & 0x7F)

	binary.BigEndian.PutUint16(buf[2:4], ff.LengthRecovery)
	binary.BigEndian.PutUint32(buf[4:8], ff.TimestampRecovery)

	binary.BigEndian.PutUint16(buf[8:10], ff.SN_base)

	binary.BigEndian.PutUint16(buf[10:12], ff.Mask)

	if ff.K1 {

		binary.BigEndian.PutUint32(buf[12:16], ff.OptionalMask1)

		// set K1		0b10100101110001001100'
		buf[10] |= 0x80
	}

	if ff.K2 {
		// set K2
		buf[12] |= 0x80
		binary.BigEndian.PutUint64(buf[16:24], ff.OptionalMask2)
	}

	return buf
}

func (ff *FecHeaderFlexibleMask) Unmarshal(buf []byte) {
	ff.R = (buf[0] >> 7 & 0x1) > 0
	ff.F = (buf[0] >> 6 & 0x1) > 0
	ff.P = (buf[0] >> 5 & 0x1) > 0
	ff.X = (buf[0] >> 4 & 0x1) > 0
	ff.CC = uint8((buf[0] & uint8(0xF)))
	ff.M = (buf[0] >> 7 & 0x1) > 0
	ff.PTRecovery = buf[1] & 0x7F
	ff.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])
	ff.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])
	ff.SN_base = binary.BigEndian.Uint16(buf[8:10])

	ff.Mask = binary.BigEndian.Uint16(buf[10:12])
	ff.Mask &= 0x7F //	unset the Most Significant bit(for k1)
	ff.K1 = (buf[10] >> 7 & 0x1) > 0
	ff.K2 = (buf[12] >> 7 & 0x1) > 0

	if ff.K1 {
		bitsM1 := buf[12:16]
		bitsM1[0] &= 0x7F
		ff.OptionalMask1 = binary.BigEndian.Uint32(bitsM1)
	}

	if ff.K2 {
		bitsM2 := buf[16:24]
		ff.OptionalMask2 = binary.BigEndian.Uint64(bitsM2) // 100101110000101100101110'

	}

}

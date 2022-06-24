package main

import (
	"encoding/binary"
)

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
	Mask1             uint16 // use 15 bits
	K2                bool
	Mask2             [3]uint32 // use 95 bits
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

func newFecHeaderFlexibleMask(R bool, F bool, P bool, X bool, CC uint8, M bool, PTRecovery uint8, LengthRecovery uint16, TimestampRecovery uint32, SN_base uint16, K1 bool, Mask1 uint16, K2 bool, Mask2 [3]uint32) FecHeader {
	return &FecHeaderFlexibleMask{
		R:                 R,
		F:                 F,
		P:                 P,
		X:                 X,
		M:                 M,
		K1:                K1,
		K2:                K2,
		CC:                CC,
		Mask1:             Mask1,
		Mask2:             Mask2,
		SN_base:           SN_base,
		PTRecovery:        PTRecovery,
		LengthRecovery:    LengthRecovery,
		TimestampRecovery: TimestampRecovery,
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

	buf[1] |= byte((ff.PTRecovery))

	binary.BigEndian.PutUint16(buf[2:4], ff.LengthRecovery)
	binary.BigEndian.PutUint32(buf[4:8], ff.TimestampRecovery)
	binary.BigEndian.PutUint16(buf[8:10], ff.SN_base)

	if ff.K1 {
		buf[10] |= (1 << 7)
		buf[10] |= byte(ff.Mask1 & 0x7F00)
		buf[11] |= byte(ff.Mask1)
	}

	if ff.K2 {
		binary.BigEndian.PutUint32(buf[12:16], ff.Mask2[0])
		buf[12] |= (1 << 7) // set K

		binary.BigEndian.PutUint32(buf[16:20], ff.Mask2[1])
		binary.BigEndian.PutUint32(buf[20:24], ff.Mask2[2])
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

	ff.K1 = (buf[10] >> 7 & 0x1) > 0
	if ff.K1 {
		bitsM1 := buf[10:12]
		bitsM1[0] &= 0x7F
		ff.Mask1 = binary.BigEndian.Uint16(bitsM1)
	}

	ff.K2 = (buf[12] >> 7 & 0x1) > 0

	if ff.K2 {
		bitsM2_1 := buf[12:16]
		bitsM2_2 := buf[16:20]
		bitsM2_3 := buf[20:24]

		bitsM2_1[0] &= 0x7F

		ff.Mask2[0] = binary.BigEndian.Uint32(bitsM2_1)
		ff.Mask2[1] = binary.BigEndian.Uint32(bitsM2_2)
		ff.Mask2[2] = binary.BigEndian.Uint32(bitsM2_3)

	}
}

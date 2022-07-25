package fech

import (
	"encoding/binary"
)

func ToFecHeaderLD(buf []byte, SN_base uint16, L, D uint8) (FecHeaderLD, []byte) {
	var fecheader FecHeaderLD
	// first 2 bits are neglected in FEC bit string and replaced by R and F
	fecheader.R = false
	fecheader.F = true
	fecheader.P = (buf[0] >> 5 & 0x1) > 0
	fecheader.X = (buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))

	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PTRecovery = buf[1] & 0x7F

	fecheader.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])

	fecheader.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])

	// Check: SN_base, L, D
	fecheader.SN_base = SN_base
	fecheader.L = L
	fecheader.D = D
	return fecheader, buf[8:]
}

func ToFecHeaderFlexibleMask(buf []byte) (FecHeaderFlexibleMask, []byte) {
	var fecheader FecHeaderFlexibleMask
	fecheader.R = false
	fecheader.F = false
	fecheader.P = (buf[0] >> 5 & 0x1) > 0
	fecheader.X = (buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))

	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PTRecovery = buf[1] & 0x7F

	fecheader.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])

	fecheader.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])

	return fecheader, buf[8:]
}

func ToFecHeaderRetransmission(buf []byte) (FecHeaderRetransmission, []byte) {
	var fecheader FecHeaderRetransmission
	fecheader.R = true
	fecheader.F = false
	fecheader.P = (buf[0] >> 5 & 0x1) > 0
	fecheader.X = (buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))

	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PayloadType = buf[1] & 0x7F

	fecheader.SeqNumber = binary.BigEndian.Uint16(buf[2:4])

	fecheader.TimeStamp = binary.BigEndian.Uint32(buf[4:8])

	return fecheader, buf[8:]
}

// function to convert the FEC bit string (type []byte) to FEC header (type FecHeaderLD)
// func ToFecHeader(buf []byte, fecvarient string) (FecHeader, []byte, error) {

// 	// using fecvarient to check
// 	// L D
// 	// flexible mask
// 	// retransmission
// 	if fecvarient == "LD" {
// 		header, body := ToFecHeaderLD(buf)
// 		return &header, body, nil
// 	}
// 	// if fecvarient == "flexible mask" {
// 	// 	header, body := ToFecHeaderFlexibleMask(buf)
// 	// 	return &header, body, nil
// 	// }
// 	if fecvarient == "retransmission" {
// 		header, body := ToFecHeaderRetransmission(buf)
// 		return &header, body, nil
// 	}
// 	return &FecHeaderLD{}, nil, errors.New("ec varient is not defined correctly.")

// }

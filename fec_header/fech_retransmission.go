package fech

import (
	"encoding/binary"
)

type FecHeaderRetransmission struct {
	R           bool
	F           bool
	P           bool
	X           bool
	CC          uint8
	M           bool
	PayloadType uint8
	SeqNumber   uint16
	TimeStamp   uint32
	SSRC        uint32
}

/*


    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |1|0|P|X|  CC   |M| Payload Type|        Sequence Number        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           Timestamp                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                              SSRC                             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   :          Retransmission "Payload" follows FEC Header          :
   :                                                               :
*/

func NewFecHeaderRetransmission(R bool, F bool, P bool, X bool, CC uint8, M bool, PayloadType uint8, SeqNumber uint16, TimeStamp uint32, SSRC uint32) FecHeader {
	return &FecHeaderRetransmission{
		R:           R,
		F:           F,
		P:           P,
		X:           X,
		M:           M,
		CC:          CC,
		SSRC:        SSRC,
		TimeStamp:   TimeStamp,
		SeqNumber:   SeqNumber,
		PayloadType: PayloadType,
	}
}

func (fh *FecHeaderRetransmission) Marshal() []byte {
	//size of fec header for retransmission = 12 bytes
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

	buf[0] |= (fh.CC & 0xF)

	if fh.M {
		buf[1] |= (1 << 7)
	}

	buf[1] |= (fh.PayloadType & 0x7F)

	binary.BigEndian.PutUint16(buf[2:4], fh.SeqNumber)
	binary.BigEndian.PutUint32(buf[4:8], fh.TimeStamp)
	binary.BigEndian.PutUint32(buf[8:12], fh.SSRC)

	return buf
}

func (fh *FecHeaderRetransmission) Unmarshal(buf []byte) {
	fh.R = (buf[0] >> 7 & 0x1) > 0
	fh.F = (buf[0] >> 6 & 0x1) > 0
	fh.P = (buf[0] >> 5 & 0x1) > 0
	fh.X = (buf[0] >> 4 & 0x1) > 0
	fh.CC = uint8((buf[0] & uint8(0xF)))

	fh.M = (buf[1] >> 7 & 0x1) > 0
	fh.PayloadType = buf[1] & 0x7F

	fh.SeqNumber = binary.BigEndian.Uint16(buf[2:4])
	fh.TimeStamp = binary.BigEndian.Uint32(buf[4:8])
	fh.SSRC = binary.BigEndian.Uint32(buf[8:12])

}

package fech

import (
	"encoding/binary"
)
// function to convert the FEC bit string (type []byte) to FEC header (type FecHeaderLD) 
func UnmarshalFec(buf []byte)(FecHeaderLD){

	// check: do we need to import FecHeaderLD or does it do it as they are in same package

	var fecheader FecHeaderLD
	// first 2 bits are neglected in FEC bit string and replaced by R and F
	fecheader.R=false
	fecheader.F=false
	fecheader.P=(buf[0] >> 5 & 0x1) > 0
	fecheader.X=(buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))
	
	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PTRecovery = buf[1] & 0x7F

	fecheader.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])
	
	fecheader.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])
	
	// Check: SN_base, L, D
	return fecheader
}
// --------------------------------------------------------------------
# flexfec
Flexible Forward Error Correction (FEC)

## Current Repository Structure
```
.
├── bitstring
│   ├── bitstring.go
│   └── fecbitstring.go
├── dump
│   ├── bitstring.go
│   ├── fecheader.go
│   ├── generateRepair.go
│   ├── recoverMissingPacket.go
│   └── generateRTP.go
├── fec_header
│   ├── fecheader.go
│   ├── fech_flexiblemask.go
│   ├── fech_fromfecbitstring.go
│   ├── fech_LD.go
│   └── fech_retransmission.go
├── go.mod
├── go.sum
├── main
│   └── test.go
├── README.md
├── testing
│   ├── bitstring_testing.go
│   ├── fecheader_testing.go
│   └── rtp_to_fech_testing.go
├── todo.md
└── util
    ├── generateRTP.go
    ├── padRTP.go
    ├── printBytes.go
    └── printRTP.go

```

# flexfec
[Flexible Forward Error Correction (FEC)](https://datatracker.ietf.org/doc/html/draft-ietf-payload-flexible-fec-scheme#section-1.1.7)

### Current Repository Structure
```
.
├── bitstring
│   ├── bitstring.go
│   ├── fecbitstring.go
│   └── getBitstrings.go
├── buffer
│   └── buffer.go
├── fec_header
│   ├── fecheader.go
│   ├── fech_flexiblemask.go
│   ├── fech_fromfecbitstring.go
│   ├── fech_LD.go
│   └── fech_retransmission.go
├── flex-fec-flow.pdf
├── go.mod
├── go.sum
├── main
│   ├── ColFec_demo.go
│   ├── FlexibleMask_demo.go
│   ├── output
│   │   ├── buffer.txt
│   │   ├── receiver.txt
│   │   └── sender.txt
│   ├── RowFec_demo.go
│   └── Two_dimension_demo.go
├── README.md
├── recover
│   ├── generate_repair.go
│   └── recover_missing_packet.go
├── testing
│   ├── 01_fecheader_testing.go
│   ├── 02_fecbitstring_testing.go
│   ├── 03_generate_repair_testing.go
│   ├── 04_buffer_testing.go
│   ├── 05_flexibleMask_testing.go
│   └── 06_flexfec_testing.go
└── util
    ├── generateRTP.go
    ├── padRTP.go
    ├── printBytes.go
    ├── printRTP.go
    └── testCaseMap.go
```


## Latest update
Working flexFEC for LD version that can recover RTP packets
with CSRC list and Extension header.

to run:

```sh
cd testing
go run .\06_flexfec_testing.go
```

Observe the debug.txt for the packets sents from sender side.

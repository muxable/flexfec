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
│   ├── repair.go
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

## Description

### .\bitstring

```go
func ToBitString(p *rtp.Packet) (out []byte)
```

takes an rtp packet and returns its bitstring(byte slice) representation as per the ieft specification.

```go
func ToFecBitString(buf [][]byte) []byte
```

takes all bitstrings of packets from the same source block and returns the FecBitstring by applying the parity code operation(XOR).

### .\fec_header

contains all the different FecHeader variants struct representation

```go
func ToFecHeader(buf []byte) (FecHeaderLD, []byte)
```

takes the FecBitstring of a source block and forms the FecHeader and RepairPayload of the repair packet construction.

### .\recover

currently coded only for 1d 1row configuartion(i.e a single source block **L>0 and D=1**)

```go
func GenerateRepair(srcBlock *[]rtp.Packet, L, D int) rtp.Packet
```

this function constructs the repair packet for the source block, basically the Encoder part

```go
func RecoverMissingPacket(srcBlock *[]rtp.Packet, repairPacket rtp.Packet) (rtp.Packet, int)
```

this function recovers the missing packet in the recieved source block, and uses the repair packet and recovers that missing packet.

### .\util

- functions that will help us test and debug: 
- generate n rtp packets(i.e a source block)
- pad packets
- print an rtp packet, or print a buffer

### .\testing

- all the functions(in other packages) testing is done under the testing package

### .\main**

- like a deployment package
- where the encoder and decoder are setup on a UDP connection

## Latest update

- recover missing packet in 1d 1row variant :  

![1d 1row fec case scenerio](https://github.com/muxable/flexfec/blob/main/dump/1d_1row_fec.png?raw=true)  

- to run:

```sh
 cd main
```

```sh
go run .\one_dim_one_row_testing.go
```

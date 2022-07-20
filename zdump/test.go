package main

import (
	"fmt"
	"net"
	"time"

	"github.com/pion/rtp"
)

const (
	listenPort = 6420
	ssrc       = 5000
	mtu        = 50
)

func sender() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		panic(err)
	}

	for i := uint16(0); ; i++ {
		time.Sleep(1 * time.Second)

		header := rtp.Header{
			Version:        2,
			SSRC:           ssrc,
			SequenceNumber: i,
		}

		payload := []byte{0x0, 0x1, 0x2}
		headerBuf, _ := header.Marshal()
		packet := append(headerBuf, payload...)

		conn.Write(packet)

		fmt.Println("sent packet :", i)
		fmt.Println(packet)
	}

}

func receiver() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", serverAddr)
	if err != nil {
		panic(err)
	}

	for {
		buffer := make([]byte, mtu)

		i, _, err := conn.ReadFrom(buffer)

		if err != nil {
			panic(err)
		}

		fmt.Println("recieved packet")
		fmt.Println("buffer : ", buffer[:i], "\n")
	}

}

func main() {
	go sender()
	receiver()
}

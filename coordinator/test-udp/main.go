package main

import (
	"fmt"
	"net"
	"os"
)

const protocol = "udp"
const address = "224.0.0.1:32100"

func main() {
	udpAddress, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		fmt.Printf("error:[%v] creating udp address\n", err)
		os.Exit(1)
	}
	udpConn, err := net.ListenMulticastUDP(protocol, nil, udpAddress)
	if err != nil {
		fmt.Printf("error:[%v] creating udp connection\n", err)
		os.Exit(1)
	}
	buf := make([]byte, 1024)
	n, addr, err := udpConn.ReadFrom(buf)
	if err != nil {
		fmt.Printf("error:[%v] creating udp connection\n", err)
		os.Exit(1)
	}
	fmt.Printf("no of bytes:%d\n", n)
	fmt.Printf("addr:%s\n", addr.String())
	fmt.Printf("data:%s\n", buf)

}

package server

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const protocol = "udp"

func buildUDPAddress() (*net.UDPAddr, error) {
	const address = "224.0.0.1:32100"
	udpAddress := os.Getenv("UDP_ADDRESS")
	if udpAddress == "" {
		udpAddress = address
	}
	return net.ResolveUDPAddr(protocol, udpAddress)
}

func UDPMultiCastListen() (*net.UDPConn, error) {
	udpAddress, err := buildUDPAddress()
	if err != nil {
		fmt.Printf("error:[%v] creating udp address\n", err)
		return nil, err
	}
	udpConn, err := net.ListenMulticastUDP(protocol, nil, udpAddress)
	if err != nil {
		fmt.Printf("error:[%v] creating udp connection\n", err)
		return nil, err
	}
	return udpConn, nil
}

var CoordinatorAddress string
var CoordinatorPort string

func Listen(conn *net.UDPConn, buf []byte) error {
	if conn == nil {
		return fmt.Errorf("conn failure")
	}
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		return err
	}
	if n < 1 {
		return fmt.Errorf("empty data received")
	}
	al := strings.Split(addr.String(), ":")
	CoordinatorAddress = al[0]
	CoordinatorPort = fmt.Sprintf("%s", buf)
	return nil

}

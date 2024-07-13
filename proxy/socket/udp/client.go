package udp

import (
	"fmt"
	"net"
	"protho/proxy/socket/protocols"
)

func Dial(destinationServer string, destinationPort int) (*net.UDPConn, *net.UDPAddr, error) {
	address, err := net.ResolveUDPAddr(protocols.UDP, fmt.Sprintf("%v:%v", destinationServer, destinationPort))
	if err != nil {
		return nil, nil, err
	}

	connection, err := net.DialUDP(protocols.UDP, nil, address)
	return connection, address, err
}

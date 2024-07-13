package tcp

import (
	"fmt"
	"net"
	"protho/proxy/socket/protocols"
)

func Dial(destinationServer string, destinationPort int) (*net.TCPConn, error) {
	address, err := net.ResolveTCPAddr(protocols.TCP, fmt.Sprintf("%v:%v", destinationServer, destinationPort))
	if err != nil {
		return nil, err
	}

	return net.DialTCP(protocols.TCP, nil, address)
}

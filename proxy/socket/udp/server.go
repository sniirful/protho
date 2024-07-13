package udp

import (
	"fmt"
	"net"
	"protho/proxy/socket/protocols"
)

type Server struct {
	address  *net.UDPAddr
	listener *net.UDPConn
}

func Listen(sourceIP string, sourcePort int) (*Server, error) {
	var err error
	server := &Server{}

	server.address, err = net.ResolveUDPAddr(protocols.UDP, fmt.Sprintf("%v:%v", sourceIP, sourcePort))
	if err != nil {
		return nil, err
	}
	server.listener, err = net.ListenUDP(protocols.UDP, server.address)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) OnPacket(
	connectionHandler func(connection *net.UDPConn, address *net.UDPAddr, bytes []byte),
	errorHandler func(err error),
	bufferSize int64,
	strictErrorChecking bool,
) error {
	buffer := make([]byte, bufferSize)
	for {
		read, address, err := server.listener.ReadFromUDP(buffer)
		if err != nil {
			if strictErrorChecking {
				return err
			} else {
				// TODO: check if it has to go on another thread
				go errorHandler(err)
				continue
			}
		}

		go connectionHandler(server.listener, address, buffer[:read])
	}
}

func (server *Server) Close() error {
	return server.listener.Close()
}

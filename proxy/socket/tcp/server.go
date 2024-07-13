package tcp

import (
	"fmt"
	"net"
	"protho/proxy/socket/protocols"
)

type Server struct {
	address  *net.TCPAddr
	listener *net.TCPListener
}

func Listen(sourceIP string, sourcePort int) (*Server, error) {
	var err error
	server := &Server{}

	server.address, err = net.ResolveTCPAddr(protocols.TCP, fmt.Sprintf("%v:%v", sourceIP, sourcePort))
	if err != nil {
		return nil, err
	}
	server.listener, err = net.ListenTCP(protocols.TCP, server.address)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) OnConnection(
	connectionHandler func(connection *net.TCPConn),
	errorHandler func(err error),
	strictErrorChecking bool,
) error {
	for {
		connection, err := server.listener.AcceptTCP()
		if err != nil {
			if strictErrorChecking {
				return err
			} else {
				// TODO: check if it has to go on another thread
				go errorHandler(err)
				continue
			}
		}

		go connectionHandler(connection)
	}
}

func (server *Server) Close() error {
	return server.listener.Close()
}

package proxy

import (
	"fmt"
	"net"
	"os"
	"protho/logs"
	"protho/proxy/socket/tcp"
)

func startTCP(
	sourceIP string,
	sourcePort int,
	destinationServer string,
	destinationPort int,
	timeout float32,
	bufferSize int64,
	verbose bool,
	strict bool,
) error {
	server, err := tcp.Listen(sourceIP, sourcePort)
	if err != nil {
		return err
	}

	return server.OnConnection(func(connection *net.TCPConn) {
		// tcp connections are pretty simple: we need both inbound
		// and outbound connection to send data to each other, at
		// the same time

		connectionAddress := connection.RemoteAddr().String()
		logs.PrintV(verbose, fmt.Sprintf("New incoming connection from %v", connectionAddress))
		defer func() {
			logs.PrintV(verbose, fmt.Sprintf("Closing inbound connection from %v", connectionAddress))
			_ = connection.Close()
		}()

		logs.PrintV(verbose, fmt.Sprintf("Connecting to %v:%v...", destinationServer, destinationPort))
		remote, err := tcp.Dial(destinationServer, destinationPort)
		if err != nil {
			logs.PrintE(fmt.Sprintf("Error connecting to %v:%v: %v", destinationServer, destinationPort, err))

			if strict {
				os.Exit(1)
			}
			return
		}
		remoteAddress := remote.RemoteAddr().String()
		defer func() {
			logs.PrintV(verbose, fmt.Sprintf("Closing outbout connection to %v", remoteAddress))
			_ = remote.Close()
		}()

		logs.PrintV(verbose, fmt.Sprintf("Forwarding: %v -> %v", connectionAddress, remoteAddress))
		logs.PrintV(verbose, fmt.Sprintf("Forwarding: %v -> %v", remoteAddress, connectionAddress))

		// what we do once we have both the inbound and
		// outbound connection is we forward each to the
		// other, and once one is done we close both
		done := make(chan bool)
		go tcp.Forward(connection, remote, timeout, bufferSize, done)
		go tcp.Forward(remote, connection, timeout, bufferSize, done)
		<-done
		close(done)
	}, func(err error) {
		logs.PrintE(fmt.Sprintf("There was an error while trying to accept an incoming connection: %v", err))
	}, strict)
}

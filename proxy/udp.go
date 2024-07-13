package proxy

import (
	"fmt"
	"net"
	"os"
	"protho/logs"
	"protho/proxy/socket/udp"
)

func startUDP(
	sourceIP string,
	sourcePort int,
	destinationServer string,
	destinationPort int,
	timeout float32,
	bufferSize int64,
	verbose bool,
	strict bool,
) error {
	server, err := udp.Listen(sourceIP, sourcePort)
	if err != nil {
		return err
	}

	// this variable is absolutely needed for associating
	// source addresses to connections towards the destination
	// server
	sessionMap := make(map[string]*net.UDPConn)
	return server.OnPacket(func(connection *net.UDPConn, address *net.UDPAddr, bytes []byte) {
		// udp connections have no real connection underneath:
		// it's just Alice sending data to Bob on the port that
		// Bob has open, but in order to do so Alice has to open
		// a port as well; this means we have to hijack Alice's
		// address and port and use that combination to store the
		// proxied connection for her, so that if another packet
		// comes from the same combination of address and port, it
		// cannot be anyone else except for Alice
		logs.PrintV(verbose, fmt.Sprintf("New incoming packet from %v", address.String()))

		// here the session key is just the address, like this:
		// 127.0.0.1:34248
		// we check if the session key is already stored and if
		// it's not, we create it
		sessionKey := address.String()
		remote, exists := sessionMap[sessionKey]
		if !exists {
			logs.PrintV(verbose, fmt.Sprintf("Connecting to %v:%v...", destinationServer, destinationPort))
			remote, _, err = udp.Dial(destinationServer, destinationPort)
			if err != nil {
				logs.PrintE(fmt.Sprintf("Error connecting to %v:%v: %v", destinationServer, destinationPort, err))

				if strict {
					os.Exit(1)
				}
				return
			}
			sessionMap[sessionKey] = remote

			// once the connection is created, we need to make sure
			// that all the data coming from that connection to Bob
			// is proxied back to Alice; since we know that all data
			// coming from that connection is going to be from Bob,
			// we can just proxy everything back to Alice
			go func() {
				// there is no need to call this function on another
				// thread and listen for data on the channel, since:
				// - we're already in another thread
				// - there's only one forward going on (see tcp.go)
				// so we can freely wait for this function call to
				// finish and then we proceed
				logs.PrintV(verbose, fmt.Sprintf("Forwarding: %v -> %v", remote.RemoteAddr().String(), address.String()))
				udp.Forward(remote, connection, address, timeout, bufferSize)

				// TODO: check if it ever stops forwarding
				logs.PrintV(verbose, fmt.Sprintf("Stopped forwarding: %v -> %v", remote.RemoteAddr().String(), address.String()))
				delete(sessionMap, sessionKey)
				logs.PrintV(verbose, fmt.Sprintf("Closing outbout connection to %v", remote.RemoteAddr().String()))
				_ = remote.Close()
			}()
		}

		logs.PrintV(verbose, fmt.Sprintf("Sending: %v -> %v", address.String(), remote.RemoteAddr().String()))
		_, _ = remote.Write(bytes)
	}, func(err error) {
		logs.PrintE(fmt.Sprintf("There was an error while trying to accept an incoming packet: %v", err))
	}, bufferSize, strict)
}

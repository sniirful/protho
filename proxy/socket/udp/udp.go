package udp

import (
	"net"
	"time"
)

// Forward reads all the data coming from sender and
// forwards it to the receiver making use of the receiver
// address; it assumes that the channel passed to the
// function is later closed by the caller
func Forward(sender, receiver *net.UDPConn, receiverAddress *net.UDPAddr, timeout float32, bufferSize int64) {
	// this is an attempt at improving performance,
	// not doing unnecessary timeout checks when not
	// needed
	if timeout > 0 {
		forwardWithTimeout(
			sender,
			receiver,
			receiverAddress,
			timeout,
			bufferSize,
		)
	} else {
		forwardWithoutTimeout(
			sender,
			receiver,
			receiverAddress,
			bufferSize,
		)
	}
}

func forwardWithTimeout(sender, receiver *net.UDPConn, receiverAddress *net.UDPAddr, timeout float32, bufferSize int64) {
	timeoutDuration := time.Duration(timeout) * time.Second
	for {
		buffer := make([]byte, bufferSize)
		// TODO: does this need error checking?
		_ = sender.SetReadDeadline(time.Now().Add(timeoutDuration))
		read, err := sender.Read(buffer)
		if err != nil {
			break
		}

		if _, err := receiver.WriteToUDP(buffer[:read], receiverAddress); err != nil {
			break
		}
	}
}

func forwardWithoutTimeout(sender, receiver *net.UDPConn, receiverAddress *net.UDPAddr, bufferSize int64) {
	for {
		buffer := make([]byte, bufferSize)
		read, err := sender.Read(buffer)
		if err != nil {
			break
		}

		if _, err := receiver.WriteToUDP(buffer[:read], receiverAddress); err != nil {
			break
		}
	}
}

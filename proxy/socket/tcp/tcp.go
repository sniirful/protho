package tcp

import (
	"net"
	"time"
)

// Forward reads all the data coming from sender and
// forwards it to the receiver; it assumes that the
// channel passed to the function is later closed by
// the caller
func Forward(sender, receiver *net.TCPConn, timeout float32, bufferSize int64, done chan bool) {
	// we want to recover as not to risk crashing
	// when sending to a closed channel
	defer func() {
		recover()
	}()

	// this is an attempt at improving performance,
	// not doing unnecessary timeout checks when not
	// needed
	if timeout > 0 {
		forwardWithTimeout(
			sender,
			receiver,
			timeout,
			bufferSize,
		)
	} else {
		forwardWithoutTimeout(
			sender,
			receiver,
			bufferSize,
		)
	}
	done <- true
}

func forwardWithTimeout(sender, receiver *net.TCPConn, timeout float32, bufferSize int64) {
	timeoutDuration := time.Duration(timeout) * time.Second
	buffer := make([]byte, bufferSize)
	for {
		// TODO: does this need error checking?
		_ = sender.SetReadDeadline(time.Now().Add(timeoutDuration))
		read, err := sender.Read(buffer)
		if err != nil {
			break
		}

		if _, err := receiver.Write(buffer[:read]); err != nil {
			break
		}
	}
}

func forwardWithoutTimeout(sender, receiver *net.TCPConn, bufferSize int64) {
	buffer := make([]byte, bufferSize)
	for {
		read, err := sender.Read(buffer)
		if err != nil {
			break
		}

		if _, err := receiver.Write(buffer[:read]); err != nil {
			break
		}
	}
}

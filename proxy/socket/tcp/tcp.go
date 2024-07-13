package tcp

import "net"

// Forward reads all the data coming from sender and
// forwards it to the receiver; it assumes that the
// channel passed to the function is later closed by
// the caller
func Forward(sender, receiver *net.TCPConn, bufferSize int64, done chan bool) {
	// we want to recover as not to risk crashing
	// when sending to a closed channel
	defer func() {
		recover()
	}()

	for {
		buffer := make([]byte, bufferSize)
		read, err := sender.Read(buffer)
		if err != nil {
			break
		}

		if _, err := receiver.Write(buffer[:read]); err != nil {
			break
		}
	}
	done <- true
}

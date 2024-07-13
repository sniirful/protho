package udp

import "net"

// Forward reads all the data coming from sender and
// forwards it to the receiver making use of the receiver
// address; it assumes that the channel passed to the
// function is later closed by the caller
func Forward(sender, receiver *net.UDPConn, receiverAddress *net.UDPAddr, bufferSize int64, done chan bool) {
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

		// TODO: this does not stop the connection when a packet
		// TODO: could not be delivered, so a timeout must be implemented
		if _, err := receiver.WriteToUDP(buffer[:read], receiverAddress); err != nil {
			break
		}
	}
	done <- true
}

package proxy

import (
	"protho/proxy/socket/protocols"
)

func Start(
	protocol string,
	sourceIP string,
	sourcePort int,
	destinationServer string,
	destinationPort int,
	bufferSize int64,
	verbose bool,
	strict bool,
) error {
	switch protocol {
	case protocols.TCP:
		return startTCP(
			sourceIP,
			sourcePort,
			destinationServer,
			destinationPort,
			bufferSize,
			verbose,
			strict,
		)
	case protocols.UDP:
		return startUDP(
			sourceIP,
			sourcePort,
			destinationServer,
			destinationPort,
			bufferSize,
			verbose,
			strict,
		)
	}

	return nil
}

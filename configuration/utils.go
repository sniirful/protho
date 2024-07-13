package configuration

import (
	"fmt"
	"protho/configuration/cutils"
	"strconv"
	"strings"
)

// getProtocol returns the protocol passed to the function,
// but only if it is included in the list of allowed protocols
func getProtocol(protocol string, allowedProtocols []string) (string, error) {
	for _, allowedProtocol := range allowedProtocols {
		if protocol == allowedProtocol {
			return protocol, nil
		}
	}

	// this creates:
	// "the protocol must be one of: 'tcp', 'udp'"
	errorString := "the protocol must be one of: "
	for i, allowedProtocol := range allowedProtocols {
		if i > 0 {
			errorString += ", "
		}
		errorString += fmt.Sprintf("'%v'", allowedProtocol)
	}
	return "", fmt.Errorf(errorString)
}

// getSplitAddress returns the address split into server
// and port; if the address does not explicitly contain
// the server, then it uses the defaultServer parameter
func getSplitAddress(address, defaultServer string) (string, int, error) {
	// this whole thing is done to support IPv6;
	// the work is as follows:
	// 1. the address is reversed (::1:8080 -> 0808:1::)
	// 2. the address is split by colons (0808:1:: -> ["0808", "1", "", ""])
	// 3. the first element of the array is now the reversed port,
	//    so we need to reverse it once more (0808 -> 8080)
	// 4. the second element onwards is the IPv6 (or IPv4) address
	//    split by colons, so we need to join it first using
	//    colons (["1", "", ""] -> 1::)
	// 5. the reversed address is then reversed once more (1:: -> ::1)
	//
	// this works well with IPv4 too:
	// 1. 127.0.0.1:8080 -> 0808:1.0.0.721
	// 2. 0808:1.0.0.721 -> ["0808", "1.0.0.721"]
	// 3. 0808 -> 8080
	// 4. ["1.0.0.721"] -> 1.0.0.721
	// 5. 1.0.0.721 -> 127.0.0.1
	reversedAddress := cutils.ReverseString(address)
	reversedSplit := strings.Split(reversedAddress, ":")
	split := []string{
		cutils.ReverseString(strings.Join(reversedSplit[1:], ":")),
		cutils.ReverseString(reversedSplit[0]),
	}
	if len(split) < 2 {
		return "", 0, fmt.Errorf("address must be in the format of [server]:port")
	}

	server := strings.Trim(split[0], " ")
	if server == "" {
		server = defaultServer
	}
	port, err := strconv.Atoi(split[1])
	if err != nil {
		return "", 0, err
	}

	return server, port, nil
}

// getBufferSize checks if the buffer size is an actual positive
// number; if it is, it returns the same buffer size passed to
// the function, and if it is not, it returns an error
func getBufferSize(bufferSize int64) (int64, error) {
	if bufferSize < 0 {
		return 0, fmt.Errorf("buffer size must be a positive number")
	}

	return bufferSize, nil
}

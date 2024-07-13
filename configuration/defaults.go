package configuration

import "protho/proxy/socket/protocols"

func getAllowedProtocols() []string {
	return []string{protocols.TCP, protocols.UDP}
}

func getDefaultServer() string {
	return "0.0.0.0"
}

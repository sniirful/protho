package main

import (
	"fmt"
	"os"
	"protho/configuration"
	"protho/logs"
	"protho/proxy"
)

func main() {
	config, err := configuration.GetConfiguration()
	if err != nil {
		logs.PrintE(fmt.Sprintf("Unable to parse configuration: %v\n", err))
		os.Exit(1)
	}

	// we repeat config so many times instead of passing it as
	// a whole because of modularity and reproducibility of the
	// proxy package
	if err = proxy.Start(
		config.Protocol,
		config.SourceIP,
		config.SourcePort,
		config.DestinationServer,
		config.DestinationPort,
		config.NoPacketTimeout,
		config.BufferSize,
		config.Verbose,
		config.Strict,
	); err != nil {
		logs.PrintE(fmt.Sprintf("Error while proxying requests: %v\n", err))
		os.Exit(1)
	}
}

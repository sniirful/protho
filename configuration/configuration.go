package configuration

import (
	"protho/configuration/cli"
)

type Configuration struct {
	SourceIP          string
	SourcePort        int
	DestinationServer string
	DestinationPort   int
	Protocol          string
	BufferSize        int64
	Verbose           bool
	Strict            bool
}

func GetConfiguration() (Configuration, error) {
	var err error
	configuration := Configuration{}
	parsedConfiguration := cli.ParseParameters()

	configuration.Protocol, err = getProtocol(parsedConfiguration.Protocol, getAllowedProtocols())
	if err != nil {
		return Configuration{}, err
	}

	configuration.SourceIP, configuration.SourcePort, err = getSplitAddress(parsedConfiguration.Source, getDefaultServer())
	if err != nil {
		return Configuration{}, err
	}
	configuration.DestinationServer, configuration.DestinationPort, err = getSplitAddress(parsedConfiguration.Destination, getDefaultServer())
	if err != nil {
		return Configuration{}, err
	}

	configuration.BufferSize = parsedConfiguration.BufferSize
	configuration.Verbose = parsedConfiguration.Verbose
	configuration.Strict = parsedConfiguration.Strict

	return configuration, nil
}

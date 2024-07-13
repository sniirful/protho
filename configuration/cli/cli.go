package cli

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	Source      string `cmd:"" name:"source" arg:"" help:"The source port (and IP) to listen on, in the format of [server]:port. If the server is not specified, it is defaulted to 0.0.0.0."`
	Destination string `cmd:"" name:"destination" arg:"" help:"The destination port (and server) to forward data to, in the format of [server]:port. If the server is not specified, it is defaulted to 0.0.0.0."`

	Protocol string `name:"protocol" short:"p" help:"The protocol to use in the communication. One of 'tcp' or 'udp'." default:"tcp"`

	NoPacketTimeout float32 `name:"timeout" short:"t" help:"How much time to wait in seconds without any packets in order to close the connection. Zero means no timeout." default:"0"`
	BufferSize      int64   `name:"buffer-size" short:"b" help:"How many bytes to proxy in one chunk at maximum." default:"65535"`
	Verbose         bool    `name:"verbose" short:"v" help:"Print more information about inbound and outbound connections." default:"false"`
	Strict          bool    `name:"strict" short:"s" help:"Be more strict about errors, do not ignore them and crash the program instead." default:"false"`
}

func ParseParameters() CLI {
	flags := CLI{}
	kong.Parse(&flags)

	return flags
}

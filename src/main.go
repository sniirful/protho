package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	Version = "0.1.0"

	CommandBufferSize      = "--buffer-size"
	CommandBufferSizeAlias = "-b"

	CommandProtocol = "--protocol"

	CommandInPort = "--in-port"

	CommandOutServer = "--out-server"
	CommandOutPort   = "--out-port"

	CommandConfigFile      = "--config-file"
	CommandConfigFileAlias = "-c"
)

var bufferSize int = 65536

var protocol string = "tcp"

var inPort string

var outServer string = "127.0.0.1"
var outPort string

var c Config

func main() {
	checkArguments()
	if inPort == "" || outPort == "" {
		printHelp()
		return
	}

	ln, err := net.Listen(protocol, fmt.Sprintf(":%v", inPort))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Listening on port %v and forwarding to %v:%v\n", inPort, outServer, outPort)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func checkArguments() {
	checkArgs(os.Args, 1, []argHandler{
		{
			args: []string{CommandBufferSize, CommandBufferSizeAlias},
			handler: func(current, next *string) {
				if next != nil {
					b, err := strconv.Atoi(*next)
					if err == nil {
						bufferSize = b
					}
				}
			},
		},
		{
			args: []string{CommandProtocol},
			handler: func(current, next *string) {
				if next != nil {
					protocol = *next
				}
			},
		},
		{
			args: []string{CommandInPort},
			handler: func(current, next *string) {
				if next != nil {
					inPort = *next
				}
			},
		},
		{
			args: []string{CommandOutServer},
			handler: func(current, next *string) {
				if next != nil {
					outServer = *next
				}
			},
		},
		{
			args: []string{CommandOutPort},
			handler: func(current, next *string) {
				if next != nil {
					outPort = *next
				}
			},
		},
		{
			args: []string{CommandConfigFile, CommandConfigFileAlias},
			handler: func(current, next *string) {
				if next != nil {
					parseConfigurationFile(*next)
				}
			},
		},
	})
}

//
// socket
//

func handleConnection(conn net.Conn) {
	remote, err := net.Dial(protocol, fmt.Sprintf("%v:%v", outServer, outPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	c := make(chan bool)
	go forward(conn, remote, c, true)
	go forward(remote, conn, c, false)
	<-c

	conn.Close()
	remote.Close()
}

func forward(sender, receiver net.Conn, c chan bool, a bool) {
	for {
		buf := make([]byte, bufferSize)
		read, err := sender.Read(buf)
		if err != nil {
			break
		}

		filtered, drop := filter(buf[:read])
		if drop {
			break
		}
		receiver.Write(filtered)
	}
	c <- true
}

func filter(in []byte) ([]byte, bool) {
	if len(c.Drop) == 0 && len(c.DropReg) == 0 && len(c.Exclude) == 0 && len(c.Replace) == 0 && len(c.ReplaceReg) == 0 {
		return in, false
	}

	out := string(in)
	for _, d := range c.Drop {
		if strings.Contains(out, d) {
			return []byte(""), true
		}
	}
	for _, d := range c.DropReg {
		m := regexp.MustCompile(d)
		if len(m.FindStringIndex(out)) > 0 {
			return []byte(""), true
		}
	}
	for _, e := range c.Exclude {
		out = strings.ReplaceAll(out, e, "")
	}
	for _, r := range c.Replace {
		out = strings.ReplaceAll(out, r.Old, r.New)
	}
	for _, r := range c.ReplaceReg {
		m := regexp.MustCompile(r.Reg)
		out = m.ReplaceAllString(out, r.New)
	}
	return []byte(out), false
}

//
// configuration
//

type Config struct {
	Drop    []string `json:"drop"`
	DropReg []string `json:"drop-reg"`
	Exclude []string `json:"exclude"`
	Replace []struct {
		Old string `json:"old"`
		New string `json:"new"`
	} `json:"replace"`
	ReplaceReg []struct {
		Reg string `json:"reg"`
		New string `json:"new"`
	} `json:"replace-reg"`
}

func parseConfigurationFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//
// arguments
//

type argHandler struct {
	args    []string
	handler func(current, next *string)
}

func checkArgs(args []string, startPos int, handlers []argHandler) {
	for i, arg := range args {
		var arg0, arg1 string
		if strings.Contains(arg, "=") {
			arr := strings.Split(arg, "=")
			arg0 = arr[0]
			arg1 = arr[1]
		}

		for _, h := range handlers {
			var equals, equals01 bool
			for _, s := range h.args {
				if arg0 == s {
					equals01 = true
				} else if arg == s {
					equals = true
				}
			}
			if !equals && !equals01 {
				continue
			}

			var current, next *string
			if equals01 {
				current = &arg0
				next = &arg1
			} else if equals {
				current = &arg
				next = nil
				if i < (len(args) - 1) {
					next = &args[i+1]
				}
			}
			h.handler(current, next)
		}
	}
}

//
// help
//

func printHelp() {
	fmt.Printf(`Protho %v (July 3rd 2022). Usage:

protho [<args>]

Following are the args for this tool. A star (*) before the
arg means it is mandatory. The round brackets below the arg
indicate the default value. The square brackets below the arg
indicate a condition for the compulsoriness.

*--in-port[=PORT]           the port this tool is going to
                            listen to to forward the packets
                            from
*--out-port[=PORT]          the port this tool is going to
                            forward the packets to
 --out-server[=SERVER]      the server this tool is going to
                            forward the packets to
   (--out-server=127.0.0.1)
 -c, --config-file[=FILE]   the configuration file that
                            specifies idk fuck it
 -b, --buffer-size[=SIZE]   the max size for received and sent
                            packets
   (-b=65536)
 --protocol[=PROTOCOL]      the protocol (tcp or udp) used for
                            communications
   (--protocol=tcp)

You can copy the following configuration as a template for the
configuration file:

{
    "_comment": "drops the tcp or udp connection whenever it finds a packet containing this",
    "drop": [
        "test-drop"
    ],
    "_comment": "drops the tcp or udp connection whenever it finds a regex match in the packet",
    "drop-reg": [
        "test[-]{1}drop-reg"
    ],
    "_comment": "removes every match from the packet",
    "exclude": [
        "test-exclude"
    ],
    "_comment": "replaces every instance of 'old' match with 'new'",
    "replace": [
        {
            "old": "test-replace-old",
            "new": "test-replace-new"
        }
    ],
    "_comment": "replaces every instance of 'reg' regex match with 'new'",
    "replace-reg": [
        {
            "reg": "test[-]{1}replace-reg",
            "new": "test-replace-reg-done"
        }
    ],
    "_comment": "note that you can remove any of 'drop', 'drop-reg', 'exclude', 'replace' and 'replace-reg' as you like"
}

`, Version)
}

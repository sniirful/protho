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
	CommandBufferSize      = "--buffer-size"
	CommandBufferSizeAlias = "-b"

	CommandInProtocol = "--in-protocol"
	CommandInPort     = "--in-port"

	CommandOutProtocol = "--out-protocol"
	CommandOutServer   = "--out-server"
	CommandOutPort     = "--out-port"

	CommandConfigFile      = "--config-file"
	CommandConfigFileAlias = "-c"
)

var bufferSize int = 65536

var inProtocol string = "tcp"
var inPort string

var outProtocol string = "tcp"
var outServer string = "127.0.0.1"
var outPort string

var c Config

func main() {
	checkArguments()
	if inPort == "" || outPort == "" {
		help()
		return
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", inPort))
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
			args: []string{CommandInProtocol},
			handler: func(current, next *string) {
				if next != nil {
					inProtocol = *next
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
			args: []string{CommandOutProtocol},
			handler: func(current, next *string) {
				if next != nil {
					outProtocol = *next
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
	remote, err := net.Dial(outProtocol, fmt.Sprintf("%v:%v", outServer, outPort))
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
		receiver.Write(filter(buf[:read]))
	}
	c <- true
}

func filter(in []byte) []byte {
	if len(c.Replace) == 0 {
		return in
	}

	out := string(in)
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
	return []byte(out)
}

//
// configuration
//

type Config struct {
	Exclude []string `json:"exclude"`
	Replace []struct {
		Old string
		New string
	} `json:"replace"`
	ReplaceReg []struct {
		Reg string
		New string
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

func help() {
	fmt.Println(`todo`)
}

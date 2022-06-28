package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	CommandInProtocol = "--in-protocol"
	CommandInPort     = "--in-port"

	CommandOutProtocol = "--out-protocol"
	CommandOutServer   = "--out-server"
	CommandOutPort     = "--out-port"
)

var inProtocol string = "tcp"
var inPort string

var outProtocol string = "tcp"
var outServer string
var outPort string

func main() {
	checkArguments()
	if inPort == "" || outServer == "" || outPort == "" {
		help()
		return
	}

	ln, _ := net.Listen("tcp", fmt.Sprintf(":%v", inPort))
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
	go forward(conn, remote, c)
	go forward(remote, conn, c)
	<-c

	conn.Close()
	remote.Close()
}

func forward(sender, receiver net.Conn, c chan bool) {
	buf := make([]byte, 1024)
	for {
		_, err := sender.Read(buf)
		if err == io.EOF {
			break
		}
		receiver.Write(buf)
	}
	c <- true
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

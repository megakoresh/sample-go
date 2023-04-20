package main

import (
	"os"
	"strings"

	"github.com/megakoresh/sample-go/send"
	"github.com/megakoresh/sample-go/server"
	"github.com/megakoresh/sample-go/util"
)

const (
	cmdSend   = "send"
	cmdServer = "server"
)

var ()

func main() {
	switch os.Args[1] {
	case cmdSend:
		util.Logger.Println("Sending")
		os.Exit(send.Send(os.Args[2:]))
	case cmdServer:
		util.Logger.Println("Launching server")
		server.Serve(os.Args[2:])
		os.Exit(0)
	}
	util.Logger.Fatalf("No recognized command found in arguments %s", strings.Join(os.Args, " "))
}

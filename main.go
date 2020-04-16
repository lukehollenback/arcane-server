package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/lukehollenback/arcane-server/handler"
	"github.com/lukehollenback/arcane-server/service/gameserverservice"
	"github.com/lukehollenback/arcane-server/util"
)

func init() {
	//
	// Initialize any packages that require explicit probing (e.g. to cause their "init()" functions
	// to fire).
	//
	handler.PkgInit()
}

func main() {
	//
	// Register a kill signal handler with the operating system so that we can gracefully shutdown if
	// necessary.
	//
	osInterrupt := make(chan os.Signal, 1)

	signal.Notify(osInterrupt, os.Interrupt)

	//
	// Load or default in the proper configuration.
	//
	tcpBindAddress := flag.String("addr", util.GetEnv("TCP_BIND_ADDRESS", "localhost"),
		"The ip address that the server should bind to for listening. Can also be specified via the "+
			"\"TCP_BIND_ADDRESS\" environment variable.")

	tcpBindPort := flag.String("tcpport", util.GetEnv("TCP_BIND_PORT", "6543"),
		"The TCP/IP port that the server should bind to for listening. Can also be specified via the "+
			"\"TCP_BIND_PORT\" environment variable.")

	flag.Parse()

	//
	// Start the Game Server Service. Note that this service will start in its own goroutine.
	//
	gameserverservice.Instance().Config(&gameserverservice.Config{
		TCPAddr:                    (*tcpBindAddress + ":" + *tcpBindPort),
		ClientHeartbeatTimeoutSecs: 60,
	})
	gameserverservice.Instance().Start()

	//
	// Block until we are shut down by the operating system.
	//
	<-osInterrupt

	gameserverservice.Instance().Stop()

	log.Print("Goodbye.")
}

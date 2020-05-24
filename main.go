package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/lukehollenback/arcane-server/handlers"
	"github.com/lukehollenback/arcane-server/services/gameserverservice"
	"github.com/lukehollenback/arcane-server/services/playerinfoservice"
	"github.com/lukehollenback/arcane-server/util"
)

func init() {
	//
	// Initialize any packages that require explicit probing (e.g. to cause their "init()" functions
	// to fire).
	//
	handlers.PkgInit()
}

func main() {
	var err error
	var ch <-chan bool

	//
	// Register a kill signal handler with the operating system so that we can gracefully shutdown if
	// necessary.
	//
	osInterrupt := make(chan os.Signal, 1)

	signal.Notify(osInterrupt, os.Interrupt)

	//
	// Load or default in the proper configuration.
	//
	tcpBindAddress := flag.String("addr", util.GetEnv(
		"TCP_BIND_ADDRESS", "localhost"),
		"The ip address that the server should bind to for listening. Can also be specified via the "+
			"\"TCP_BIND_ADDRESS\" environment variable.",
	)

	tcpBindPort := flag.String(
		"tcpport", util.GetEnv("TCP_BIND_PORT", "6543"),
		"The TCP/IP port that the server should bind to for listening. Can also be specified via the "+
			"\"TCP_BIND_PORT\" environment variable.",
		)

	flag.Parse()

	//
	// Start the Player Info Service.
	//
	ch, err = playerinfoservice.Instance().Start()
	if err != nil {
		log.Fatalf("Failed to start the Player Info Service. (Error: %s)", err)
	}

	<-ch

	//
	// Start the Game Server Service.
	//
	gameserverservice.Instance().Config(&gameserverservice.Config{
		TCPAddr:                    *tcpBindAddress + ":" + *tcpBindPort,
		ClientHeartbeatTimeoutSecs: 60,
	})
	ch, err = gameserverservice.Instance().Start()
	if err != nil {
		log.Fatalf("Failed to start the Game Server Service. (Error: %s)", err)
	}

	<-ch

	//
	// Log some debug info.
	//
	log.Print("All services are now online.")

	//
	// Block until we are shut down by the operating system.
	//
	<-osInterrupt

	log.Print("An operating system interrupt has been received. Shutting down all services...")

	//
	// Shut down the Game Server Service.
	//
	ch, err = gameserverservice.Instance().Stop()
	if err != nil {
		log.Fatalf("Failed to stop the Game Server Service. (Error: %s)", err)
	}

	<-ch

	//
	// Shut down the Player Info Service.
	//
	ch, err = playerinfoservice.Instance().Stop()
	if err != nil {
		log.Fatalf("Failed to stop the Player Info Service. (Error: %s)", err)
	}

	<-ch

	//
	// Wrap everything up.
	//
	log.Print("Goodbye.")
}

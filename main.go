package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/archit120/tcptun/client"
	"github.com/archit120/tcptun/server"
)

func main() {

	isClient := flag.Bool("client", false, "Start as a vpn client")
	isServer := flag.Bool("server", false, "Start as a vpn server")
	port := flag.Int("l", 443, "Server listening port")
	serverip := flag.String("s", "", "Server listening port")
	debug := flag.Bool("debug", false, "show debug info")
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if *isClient && *isServer {
		logrus.Fatal("Both client and server options specified")
	} else if !(*isClient || *isServer) {
		logrus.Fatal("Neither client not server option specified")
	} else if *isClient {
		logrus.Info("Starting as client\n")
		client.StartClient(*serverip)
	} else {
		logrus.Info("Starting as server\n")
		server.StartServer(*port)
	}
}

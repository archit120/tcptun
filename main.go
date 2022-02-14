package main

import (
	"flag"
	"log"

	"github.com/archit120/tcptun/client"
	"github.com/archit120/tcptun/server"
)

func main() {

	isClient := flag.Bool("client", false, "Start as a vpn client")
	isServer := flag.Bool("server", false, "Start as a vpn server")
	port := flag.Int("l", 443, "Server listening port")
	serverip := flag.String("s", "", "Server listening port")

	flag.Parse()

	if *isClient && *isServer {
		log.Fatal("Both client and server options specified")
	} else if !(*isClient || *isServer) {
		log.Fatal("Neither client not server option specified")
	} else if *isClient {
		log.Printf("Starting as client\n")
		client.StartClient(*serverip)
	} else {
		log.Printf("Starting as server\n")
		server.StartServer(*port)
	}
}

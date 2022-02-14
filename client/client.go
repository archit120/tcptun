package client

import (
	"log"
	"net"

	"github.com/songgao/water"
)

func StartClient(serverIP string) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	defer ifce.Close()
	// Create TCP connection to server

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	// execute ip commands to activate the interface and setup routes

	packet := make([]byte, 2000)
	var packet2 [4]byte
	packet2[0] = 1
	for {
		n, err := ifce.Read(packet)
		log.Printf("Received packet of size %d sending to conn.\n", n)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(packet[:n])
	}

}

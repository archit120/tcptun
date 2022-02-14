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

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			n, err = ifce.Write(buf[:n])
			if err != nil {
				break
			}
		}
	}()

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		log.Printf("Received packet of size %d sending to conn.\n", n)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(packet[:n])
	}

}

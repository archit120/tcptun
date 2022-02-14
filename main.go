package main

import (
	"log"

	"github.com/songgao/water"
)

func main() {
	config := water.Config{
		DeviceType: water.TUN,
	}
	// config.Name = "O_O"

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	packet := make([]byte, 2000)

	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Dst: %x\n", packet[:n])
	}
}

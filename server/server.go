package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/songgao/water"
)

func StartServer(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	config := water.Config{
		DeviceType: water.TUN,
	}
	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	defer ifce.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			buf := make([]byte, 1500)
			for {
				n, err := ifce.Read(buf)
				if err != nil {
					break
				}
				n, err = conn.Write(buf[:n])
				if err != nil {
					break
				}
			}
		}()

		buf := make([]byte, 1500)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			receivedPacket := buf[:n]
			log.Printf("%x\n", receivedPacket)
			log.Printf("Recieved packet of %d from client\n", n)
			log.Printf("Dumping on this interface\n")
			ifce.Write(buf[:n])
		}
		conn.Close()
	}



}

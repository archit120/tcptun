package client

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/archit120/tcptun/common"
)

func StartClient(serverIP string) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	ifce, err := water.New(config)
	if err != nil {
		logrus.Fatal(err)
	}

	defer ifce.Close()
	// Create TCP connection to server

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		logrus.Fatal(err)
	}

	defer conn.Close()
	// execute ip commands to activate the interface and setup routes

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := common.ReadPackedPacket(conn, buf)
			if err != nil {
				logrus.Error("Error in client connection read")
				logrus.Error(err)
				break
			}
			n, err = ifce.Write(buf[:n])
			if err != nil {
				logrus.Error("Error in client interface write")
				logrus.Error(err)
				break
			}
		}
	}()

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		logrus.Debug("Received packet of size %d sending to conn.\n", n)
		if err != nil {
			logrus.Fatal(err)
		}
		common.WritePackedPacket(conn, packet[:n])
	}

}

package client

import (
	"bufio"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/archit120/tcptun/common"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func StartClient(serverIP string) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	ifce, err := water.New(config)
	if err != nil {
		logrus.Fatal(err)
	}

	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        exec.Command("/bin/sh", "./scripts/client_cleanup.sh", strings.Split(serverIP, ":")[0]).Output()
        os.Exit(0)
    }()


	defer ifce.Close()
	// Create TCP connection to server

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		logrus.Fatal(err)
	}

	defer conn.Close()
	// execute ip commands to activate the interface and setup routes
	logrus.Info("Running config script")
	cmd, err := exec.Command("/bin/sh", "./scripts/client.sh", "192.168.200.2/24", ifce.Name(), strings.Split(serverIP, ":")[0]).Output()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(cmd)

    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := common.ReadPackedPacket(reader, buf)
			
			if err != nil {
				logrus.Error("Error in client connection read")
				logrus.Error(err)
				break
			}
			logrus.Debug("Received packet from server of size ", n)
			n, err = ifce.Write(buf[:n])
			if err != nil {
				logrus.Error("Error in client interface write")
				logrus.Error(err)
				break
			}
		}
		conn.Close()
	}()

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		logrus.Debug("Received packet of size %d sending to conn.\n", n)
		if err != nil {
			logrus.Fatal(err)
		}
		common.WritePackedPacket(writer, packet[:n])
	}

}

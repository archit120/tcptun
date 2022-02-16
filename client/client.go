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
		DeviceType: water.TAP,
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
	c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        exec.Command("./scripts/client_cleanup.sh", strings.Split(serverIP, ":")[0], ifce.Name()).Output()
        conn.Close() // this will casue main to end aswell
		ifce.Close()
    }()

	defer conn.Close()
	// execute ip commands to activate the interface and setup routes
	logrus.Info("Running config script", ifce.Name())
	cmd, err := exec.Command("./scripts/client_p1.sh", ifce.Name()).Output()
	if err != nil {
		logrus.Info(string(cmd))
		logrus.Fatal(err)
	}
	logrus.Info("Script 1 done")
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
	
	go func() {
		cmd, err:= exec.Command("./scripts/client_p2.sh", ifce.Name(), strings.Split(serverIP, ":")[0]).Output()
		if err != nil {
			logrus.Info(string(cmd))
			logrus.Fatal(err)
		}
		logrus.Info("Script 2 done")
	}()

	packet := make([]byte, 1500)
	for {
		n, err := ifce.Read(packet)
		logrus.Debug("Received packet of size %d sending to conn.\n", n)
		if err != nil {
			logrus.Fatal(err)
		}
		common.WritePackedPacket(writer, packet[:n])
	}

}

package client

import (
	"bufio"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/archit120/tcptun/common"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func StartClient(serverIP string) {
		// Create TCP connection to server
	logrus.Info("Running initial script")
	cmd, err := exec.Command("./scripts/client_p0.sh", strings.Split(serverIP, ":")[0]).Output()
	if err != nil {
		logrus.Info(string(cmd))
		logrus.Fatal(err)
	}

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		logrus.Fatal(err)
	}
	
	
	config := water.Config{
		DeviceType: water.TAP,
	}
	ifce, err := water.New(config)
	if err != nil {
		conn.Close()
		logrus.Fatal(err)
	}

	defer ifce.Close()
	cleanup := func() {
		conn.Close()
		ifce.Close()
		exec.Command("./scripts/client_cleanup.sh", strings.Split(serverIP, ":")[0], ifce.Name()).Output()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	defer conn.Close()
	// execute ip commands to activate the interface and setup routes
	if runtime.GOOS != "windows" {
		logrus.Info("Running config script", ifce.Name())
		cmd, err := exec.Command("./scripts/client_p1.sh", ifce.Name(), strings.Split(serverIP, ":")[0]).Output()
		if err != nil {
			logrus.Info(string(cmd))
			cleanup()
			logrus.Fatal(err)
		}
	}
	logrus.Info("Script 1 done")
	reader := bufio.NewReader(conn)
	writer := io.Writer(conn)

	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := common.ReadPackedPacket(reader, buf)

			if err != nil {
				logrus.Error("Error in client connection read")
				logrus.Error(err)
				cleanup()
			}
			logrus.Debug("Received packet from server of size ", n)
			n, err = ifce.Write(buf[:n])
			if err != nil {
				logrus.Error("Error in client interface write")
				logrus.Error(err)

				cleanup()
			}
		}
	}()

	go func() {
		var err error
		var cmd []byte
		if runtime.GOOS != "windows" {
			cmd, err = exec.Command("./scripts/client_p2.sh", ifce.Name(), strings.Split(serverIP, ":")[0]).Output()
		} else {
			cmd, err = exec.Command("./scripts/client_p2.sh", ifce.Name(), strings.Split(serverIP, ":")[0]).Output()
		}
		if err != nil {
			logrus.Info(string(cmd))
			cleanup()
			logrus.Fatal(err)
		}

		logrus.Info("Script 2 done")
	}()

	packet := make([]byte, 1500)
	for {
		n, err := ifce.Read(packet)
		logrus.Debug("Received packet of size %d sending to conn.\n", n)
		if err != nil {
			cleanup()
			logrus.Fatal(err)
		}
		common.WritePackedPacket(writer, packet[:n])
	}

}

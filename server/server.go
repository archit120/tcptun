package server

import (
	"bufio"
	"io"
	"net"
	"os/exec"
	"strconv"

	"github.com/archit120/tcptun/common"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func StartServer(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		logrus.Fatal(err)
	}

	defer listener.Close()

	config := water.Config{
		DeviceType: water.TAP,
	}
	ifce, err := water.New(config)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Interface name is " +ifce.Name())
	logrus.Info("Running config script")
	cmd, err := exec.Command("/bin/sh", "./scripts/server.sh", "192.168.200.1/24", ifce.Name()).Output()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(cmd)
	defer ifce.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Fatal(err)
		}
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
	
		go func() {
			buf := make([]byte, 1500)

			for {
				n, err := ifce.Read(buf)
				if err != nil {
					logrus.Error("Error in server interface read")
					logrus.Error(err)
					break
				}
				n, err = common.WritePackedPacket(writer, buf[:n])
				if err != nil {
					logrus.Error("Error in server connection write")
					logrus.Error(err)
					break
				}
			}

			conn.Close()
		}()

		buf := make([]byte, 1500)
		for {
			n, err := common.ReadPackedPacket(reader, buf)
			if err != nil {
				if err != io.EOF {
					logrus.Error(err)
				}
				break
			}
			logrus.Debug("Recieved packet of %d from client\n", n)
			logrus.Debug("Dumping on this interface\n")
			ifce.Write(buf[:n])
		}
		conn.Close()
	}



}

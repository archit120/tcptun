package server

import (
	"io"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
	"github.com/archit120/tcptun/common"

)

func StartServer(port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		logrus.Fatal(err)
	}

	defer listener.Close()

	config := water.Config{
		DeviceType: water.TUN,
	}
	ifce, err := water.New(config)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Interface name is " +ifce.Name())
	defer ifce.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Fatal(err)
		}

		go func() {
			buf := make([]byte, 1500)
			for {
				n, err := ifce.Read(buf)
				if err != nil {
					logrus.Error("Error in server interface read")
					break
				}
				n, err = common.WritePackedPacket(conn, buf[:n])
				if err != nil {
					logrus.Error("Error in server connection write")

					break
				}
			}
		}()

		buf := make([]byte, 1500)
		for {
			n, err := common.ReadPackedPacket(conn, buf)
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

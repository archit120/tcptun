package common

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/sirupsen/logrus"
)

func ReadPackedPacket(reader *bufio.Reader, buffer []byte) (int, error) {
	n, err := io.ReadFull(reader, buffer[:2])
	if err!= nil {
		return 0, err
	}
	size := binary.BigEndian.Uint16(buffer[:2])
	if int(size) > len(buffer) {
		logrus.Error(size)
		return 0, errors.New("Insufficient buffer size ")
	}
	logrus.Debug(size)
	n, err = io.ReadFull(reader, buffer[:int(size)])
	if err!= nil{
		return 0, err
	}
	if n!= int(size) {
		return 0, errors.New("Incomplete data")
	}
	return n, err
}

func WritePackedPacket(writer io.Writer, buffer []byte) (int, error) {
	sizebuff := make([]byte, 2)
	binary.BigEndian.PutUint16(sizebuff, uint16(len(buffer)))
	if len(buffer) > 1500 {
		logrus.Error(len(buffer))
	}
	// log.Printf("Size is %d %x", len(buffer), sizebuff)
	n, err := writer.Write(sizebuff)
	if err!= nil {
		return 0, err
	} else if n != 2 {
		return 0, errors.New("Didnt send enough bytes")
	}
	n, err = writer.Write(buffer)
	if err!= nil {
		return 0, err
	} else if n != len(buffer) {
		return 0, errors.New("Didnt send enough bytes")
	}
	return n, err
}
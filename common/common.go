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
	logrus.Debug(size)
	n, err = io.ReadAtLeast(reader, buffer, int(size))
	if err!= nil{
		return 0, err
	}
	return n, err
}

func WritePackedPacket(writer *bufio.Writer, buffer []byte) (int, error) {
	sizebuff := make([]byte, 2)
	binary.BigEndian.PutUint16(sizebuff, uint16(len(buffer)))
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
	err = writer.Flush()
	return n, err
}
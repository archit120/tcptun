package common

import (
	"encoding/binary"
	"errors"
	"io"

	// "log"
	"net"
)

func ReadPackedPacket(reader io.Reader, buffer []byte) (int, error) {
	n, err := io.ReadAtLeast(reader, buffer, 2)
	if err!= nil && err!=io.EOF {
		return 0, err
	}
	size := binary.BigEndian.Uint16(buffer[:2])
	n, err = io.ReadAtLeast(reader, buffer, int(size))
	if err!= nil && err!=io.EOF{
		return 0, err
	}
	return n, err
}

func WritePackedPacket(conn net.Conn, buffer []byte) (int, error) {
	sizebuff := make([]byte, 2)
	binary.BigEndian.PutUint16(sizebuff, uint16(len(buffer)))
	// log.Printf("Size is %d %x", len(buffer), sizebuff)
	n, err := conn.Write(sizebuff)
	if err!= nil {
		return 0, err
	} else if n != 2 {
		return 0, errors.New("Didnt send enough bytes")
	}
	n, err = conn.Write(buffer)
	if err!= nil {
		return 0, err
	} else if n != len(buffer) {
		return 0, errors.New("Didnt send enough bytes")
	}
	return n, err
}
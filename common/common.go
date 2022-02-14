package common

import (
	"encoding/binary"
	"errors"
	// "log"
	"net"
)

func ReadPackedPacket(conn net.Conn, buffer []byte) (int, error) {
	n, err := conn.Read(buffer[:2])
	if err!= nil {
		return 0, err
	} else if n != 2 {
		return 0, errors.New("Didnt return enough bytes")
	}
	size := binary.BigEndian.Uint16(buffer[:2])
	n, err = conn.Read(buffer[:size])
	if err!= nil {
		return 0, err
	} else if n != int(size) {
		return 0, errors.New("Didnt return enough bytes")
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
		return 0, errors.New("Didnt return enough bytes")
	}
	return n, err
}
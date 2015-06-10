package fressian

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
)

type UUID struct {
	Msb, Lsb uint64
}

func NewUUID() UUID {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal("NewUUID (rand.Read): ", err)
	}
	msb := binary.BigEndian.Uint64(buf[0:8])
	lsb := binary.BigEndian.Uint64(buf[8:])
	return UUID{msb, lsb}
}

func NewUUIDFromBytes(buf []byte) UUID {
	msb := binary.BigEndian.Uint64(buf[0:8])
	lsb := binary.BigEndian.Uint64(buf[8:])
	return UUID{msb, lsb}
}

func (u UUID) Bytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[0:8], u.Msb)
	binary.BigEndian.PutUint64(buf[8:], u.Lsb)
	return buf
}

func (u UUID) String() string {
	buf := u.Bytes()
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

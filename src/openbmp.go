package go-openbmp

import (
	"bytes"
	"fmt"
	"io"
	"encoding/binary"
)

type OpenBMPHeader struct {
        Magic [4]byte
        MajorVersion uint8
        MinorVersion uint8
        HeaderLength uint16
        MessageLength uint32
	Flags uint8
	ObjType uint8
	CollectorTimestamp uint32
	CollectorTimestampMs uint32
	CollectorHash [16]byte
	CollectorAdminIDLength uint16
	CollectorAdminID []byte
	RouterHash [16]byte
	RouterIP [16]byte
	RouterGroupLength uint16
	RouterGroup []byte
	RowCount uint32
}

func ParseHeader(data []byte) OpenBMPHeader {
	message := OpenBMPHeader{}
	b1 := bytes.NewReader(data)
	binary.Read(b1, binary.BigEndian, &message.Magic)
	binary.Read(b1, binary.BigEndian, &message.MajorVersion)
	binary.Read(b1, binary.BigEndian, &message.MinorVersion)
	binary.Read(b1, binary.BigEndian, &message.HeaderLength)
	binary.Read(b1, binary.BigEndian, &message.MessageLength)
	binary.Read(b1, binary.BigEndian, &message.Flags)
	binary.Read(b1, binary.BigEndian, &message.ObjType)
	binary.Read(b1, binary.BigEndian, &message.CollectorTimestamp)
	binary.Read(b1, binary.BigEndian, &message.CollectorTimestampMs)
	binary.Read(b1, binary.BigEndian, &message.CollectorHash)
	binary.Read(b1, binary.BigEndian, &message.CollectorAdminIDLength)

	message.CollectorAdminID = make([]byte, message.CollectorAdminIDLength)
	_, err := io.ReadFull(b1, message.CollectorAdminID)
	if err != nil {
		// handle error
		fmt.Println("Error parsing!")
	}

	binary.Read(b1, binary.BigEndian, &message.RouterHash)
	binary.Read(b1, binary.BigEndian, &message.RouterIP)


	return message
}



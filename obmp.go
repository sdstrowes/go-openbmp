package obmp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
	"net/netip"
	"encoding/binary"
)

type OpenBMPHeader struct {
        MajorVersion uint8
        MinorVersion uint8
        HeaderLength uint16
        MessageLength uint32
	Flags uint8
	ObjType uint8

	CollectorTimestamp time.Time

	CollectorHash string
	CollectorAdminID []byte

	RouterHash string
	RouterIP netip.Addr

	RouterGroup []byte

	RowCount uint32

	BMPMessage []byte
}

func magicNumberIsGood(input *bytes.Reader) bool {
	var tmp [4]byte
	binary.Read(input, binary.BigEndian, &tmp)
	magic := [4]byte{0x4F, 0x42, 0x4D, 0x50}

	if magic != tmp {
		return false
	}

	return true
}

func parseTimestamp(input *bytes.Reader) (time.Time, error) {
	var collectorTimestamp uint32
	var collectorTimestampUs uint32
	binary.Read(input, binary.BigEndian, &collectorTimestamp)
	binary.Read(input, binary.BigEndian, &collectorTimestampUs)

	t := time.Unix(int64(collectorTimestamp), int64(collectorTimestampUs))

	return t, nil
}

func parse16byteHash(input *bytes.Reader, hash *string) {
	*hash = "0x"
	var tmpHash [16]byte
	binary.Read(input, binary.BigEndian, &tmpHash)
	for i := 0; i < len(tmpHash); i++ {
		*hash = *hash + fmt.Sprintf("%02x", tmpHash[i])
	}
}

func ParseHeader(data []byte) (*OpenBMPHeader, error) {
	message := OpenBMPHeader{}
	input := bytes.NewReader(data)

	if !magicNumberIsGood(input) {
		return nil, errors.New("Bad magic number")
	}

	binary.Read(input, binary.BigEndian, &message.MajorVersion)
	binary.Read(input, binary.BigEndian, &message.MinorVersion)
	binary.Read(input, binary.BigEndian, &message.HeaderLength)
	binary.Read(input, binary.BigEndian, &message.MessageLength)
	binary.Read(input, binary.BigEndian, &message.Flags)

	isRouterMessage := message.Flags & 0x80
	isRouterIPv4    := message.Flags & 0x40
	if isRouterMessage == 0 {
		return nil, errors.New("Not a router message")
	}

	binary.Read(input, binary.BigEndian, &message.ObjType)

	// parse timestamp
	message.CollectorTimestamp, _ = parseTimestamp(input)

	// parse collector has
	parse16byteHash(input, &message.CollectorHash)

	// parse collector admin ID
	var tmpCollectorAdminIDLength uint16
	binary.Read(input, binary.BigEndian, &tmpCollectorAdminIDLength)
	message.CollectorAdminID = make([]byte, tmpCollectorAdminIDLength)
	_, err := io.ReadFull(input, message.CollectorAdminID)
	if err != nil {
		fmt.Println("Error parsing!")
		return nil, err
	}

	// parse router hash
	parse16byteHash(input, &message.RouterHash)

	// parse collector IP address
	var tmpIPData [16]byte
	binary.Read(input, binary.BigEndian, &tmpIPData)
	if isRouterIPv4 == 0 {
		var tmp [4]byte
		copy(tmp[:], tmpIPData[0:4])
		message.RouterIP = netip.AddrFrom4(tmp)
	} else {
		message.RouterIP = netip.AddrFrom16(tmpIPData)
	}

	// parse router group info
	var routerGroupLength uint16
	binary.Read(input, binary.BigEndian, &routerGroupLength)
	if routerGroupLength == 0 {
		message.RouterGroup = []byte{}
	} else {
		message.RouterGroup = make([]byte, routerGroupLength)
		_, err := io.ReadFull(input, message.RouterGroup)
		if err != nil {
			fmt.Println("Error parsing!")
			return nil, err
		}
	}

	// parse rows
	binary.Read(input, binary.BigEndian, &message.RowCount)

	// read remaining bytes
	message.BMPMessage = make([]byte, input.Len())
	binary.Read(input, binary.BigEndian, &message.BMPMessage)

	return &message, nil
}


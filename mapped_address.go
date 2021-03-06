package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const (
	IPv4 = 0x01
)

type MappedAddress struct {
	Family  uint8
	Port    uint16
	Address uint32
}

func NewMappedAddress(family uint8, port uint16, address uint32) *MappedAddress {
	return &MappedAddress{
		Family:  family,
		Port:    port,
		Address: address,
	}
}

func (mappedAddress *MappedAddress) IPAddress() net.IP {
	address := make([]byte, 4)
	binary.BigEndian.PutUint32(address, mappedAddress.Address)
	ipAddress := net.IPv4(address[0], address[1], address[2], address[3])
	return ipAddress
}

func (mappedAddress *MappedAddress) Length() uint16 {
	switch mappedAddress.Family {
	case IPv4:
		return 5
	}

	return 0
}

func (mappedAddress *MappedAddress) Serialize() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, mappedAddress.Family)
	binary.Write(buffer, binary.BigEndian, mappedAddress.Port)
	binary.Write(buffer, binary.BigEndian, mappedAddress.Address)
	return buffer.Bytes()
}

func ParseMappedAddress(rawMappedAddress []byte) (*MappedAddress, error) {
	buffer := bytes.NewBuffer(rawMappedAddress)
	mappedAddress := &MappedAddress{}

	var alignment uint8
	binary.Read(buffer, binary.BigEndian, &alignment)
	if alignment != 0x00 {
		return nil, errors.New("Mapped address padding is not empty")
	}

	binary.Read(buffer, binary.BigEndian, &mappedAddress.Family)
	if mappedAddress.Family != IPv4 {
		return nil, errors.New("Mapped address family is invalid")
	}

	binary.Read(buffer, binary.BigEndian, &mappedAddress.Port)
	binary.Read(buffer, binary.BigEndian, &mappedAddress.Address)

	return mappedAddress, nil
}

func ParseXORMappedAddress(rawMappedAddress []byte, cookie uint32) (*MappedAddress, error) {
	buffer := bytes.NewBuffer(rawMappedAddress)
	mappedAddress := &MappedAddress{}

	var alignment uint8
	binary.Read(buffer, binary.BigEndian, &alignment)
	if alignment != 0x00 {
		return nil, errors.New("Mapped address padding is not empty")
	}

	binary.Read(buffer, binary.BigEndian, &mappedAddress.Family)
	if mappedAddress.Family != IPv4 {
		return nil, errors.New("Mapped address family is invalid")
	}

	binary.Read(buffer, binary.BigEndian, &mappedAddress.Port)
	binary.Read(buffer, binary.BigEndian, &mappedAddress.Address)

	cookieByte := make([]byte, 4)
	binary.BigEndian.PutUint32(cookieByte, cookie)

	cookieBuffer := bytes.NewBuffer(cookieByte) 

	var cookie16 uint16
	binary.Read(cookieBuffer, binary.BigEndian, &cookie16)

	mappedAddress.Port = mappedAddress.Port ^ cookie16
	mappedAddress.Address = mappedAddress.Address ^ cookie

	return mappedAddress, nil
}

func (mappedAddress *MappedAddress) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Family: %d\n", mappedAddress.Family))
	buffer.WriteString(fmt.Sprintf("Port: %d\n", mappedAddress.Port))
	buffer.WriteString(fmt.Sprintf("Address: %d\n", mappedAddress.Address))
	return buffer.String()
}

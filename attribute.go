package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	MappedAddressAttribute = 0x0001
	XORRelayAddressAttribute = 0x0016
	XORMappedAddressAttribute = 0x0020
)

type AttributeValue interface {
	Serialize() []byte
	Length() uint16
	String() string
}

type Attribute struct {
	Type   uint16
	Length uint16
	Value  AttributeValue
}

func NewAttribute(attributeType uint16, value AttributeValue) *Attribute {
	return &Attribute{
		Type:   attributeType,
		Length: value.Length(),
		Value:  value,
	}
}

func (attribute *Attribute) Serialize() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, attribute.Type)
	binary.Write(buffer, binary.BigEndian, attribute.Length)
	bytes := buffer.Bytes()

	bytes = append(bytes, attribute.Value.Serialize()...)

	return bytes
}

func ParseAttributes(rawAttributes []byte, header *Header) ([]*Attribute, error) {
	buffer := bytes.NewBuffer(rawAttributes)
	attributes := []*Attribute{}

	for buffer.Len() > 0 {
		attribute := &Attribute{}
		binary.Read(buffer, binary.BigEndian, &attribute.Type)
		binary.Read(buffer, binary.BigEndian, &attribute.Length)

		rawValue := make([]byte, attribute.Length)
		binary.Read(buffer, binary.BigEndian, &rawValue)
		value, err := ParseAttributeValue(rawValue, attribute.Type, header)
		if err != nil {
			continue
		}

		attribute.Value = value

		attributes = append(attributes, attribute)
	}

	return attributes, nil
}

func ParseAttributeValue(rawValue []byte, attributeType uint16, header *Header) (AttributeValue, error) {
	switch attributeType {
	case MappedAddressAttribute:
		return ParseMappedAddress(rawValue)
	case XORMappedAddressAttribute:
		return ParseXORMappedAddress(rawValue, header.MagicCookie)
	case XORRelayAddressAttribute:
		return ParseXORMappedAddress(rawValue, header.MagicCookie)
	}

	return nil, errors.New("Attribute type is invalid")
}

func (attribute *Attribute) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Type: %d\n", attribute.Type))
	buffer.WriteString(fmt.Sprintf("Length: %d\n", attribute.Length))
	buffer.WriteString("Value:\n")
	buffer.WriteString(attribute.Value.String())
	return buffer.String()
}

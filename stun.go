package stun

import (
	"errors"
	"net"
	"time"
)

const (
	GoogleStunServer           = "stun.l.google.com:19302"
	MaxResponseLength          = 548
	RequestTimeoutMilliseconds = 500
)

var StunServer = GoogleStunServer

func RequestPublicIPAddress() (net.IP, error) {
	responseMessage, err := Request(RequestClass, BindingMethod)
	if err != nil {
		return nil, err
	}

	attributeValue := responseMessage.Attributes[0].Value
	mappedAddress, ok := attributeValue.(*MappedAddress)
	if !ok {
		return nil, errors.New("Attribute was expected to be of type MappedAddress")
	}

	return mappedAddress.IPAddress(), nil
}

func RequestAllocate() (net.IP, uint16, error) {
	responseMessage, err := Request(RequestClass, AllocateMethod)
	if err != nil {
		return nil, 0, err
	}

	attributeValue := responseMessage.Attributes[0].Value
	mappedAddress, ok := attributeValue.(*MappedAddress)
	if !ok {
		return nil, 0, errors.New("Attribute was expected to be of type MappedAddress")
	}

	return mappedAddress.IPAddress(), mappedAddress.Port, nil
}

func Request(class uint16, method uint16) (*Message, error) {
	message := &Message{
		Header:     NewHeader(class, method),
		Attributes: []*Attribute{},
	}

	return RequestMessage(message)
}

func RequestMessage(request *Message) (*Message, error) {
	connection, err := net.DialTimeout("udp", StunServer, RequestTimeout())
	if err != nil {
		return nil, err
	}

	defer connection.Close()

	_, err = connection.Write(request.Serialize())
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, MaxResponseLength)
	readBytes, err := connection.Read(buffer)
	if err != nil {
		return nil, err
	}

	buffer = buffer[0:readBytes]

	return ParseMessage(buffer)
}

func RequestTimeout() time.Duration {
	return time.Duration(RequestTimeoutMilliseconds) * time.Millisecond
}

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

func RequestPublicIPAddress() (net.IP, error) {
	return RequestPublicIPAddressWithServer(GoogleStunServer)
}

func RequestPublicIPAddressWithServer(server string) (net.IP, error) {
	responseMessage, err := Request(RequestClass, BindingMethod, GoogleStunServer)
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
	return RequestAllocateWithServer(GoogleStunServer)
}

func RequestAllocateWithServer(server string) (net.IP, uint16, error) {
	responseMessage, err := Request(RequestClass, AllocateMethod, server)
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

func Request(class uint16, method uint16, server string) (*Message, error) {
	message := &Message{
		Header:     NewHeader(class, method),
		Attributes: []*Attribute{},
	}

	return RequestMessage(message, server)
}

func RequestMessage(request *Message, server string) (*Message, error) {
	connection, err := net.DialTimeout("udp", server, RequestTimeout())
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

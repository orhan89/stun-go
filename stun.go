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

func RequestAllocate(session_id int) (net.IP, uint16, error) {
	return RequestAllocateWithServer(GoogleStunServer, session_id)
}

func RequestAllocateWithServer(server string, session_id int) (net.IP, uint16, error) {
	var sessionID SessionID
	sessionID = SessionID(session_id)
	session_attribute := NewAttribute(SessionIDAttribute, &sessionID)

	responseMessage, err := Request(RequestClass, AllocateMethod, server, session_attribute)
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

func RequestRefresh(lifetime int) (error) {
	return RequestRefreshWithServer(GoogleStunServer, lifetime)
}

func RequestRefreshWithServer(server string, lifetime int) (error) {
	var lifetime_value Lifetime = Lifetime(lifetime)
	lifetime_attribute := NewAttribute(LifetimeAttribute, &lifetime_value)

	_, err := Request(RequestClass, RefreshMethod, server, lifetime_attribute)
	if err != nil {
		return err
	}
	return nil
}

func Request(class uint16, method uint16, server string, attributes ...*Attribute) (*Message, error) {
	if attributes == nil {
		attributes = []*Attribute{}
	}

	message := &Message{
		Header:     NewHeader(class, method),
		Attributes: attributes,
		Padding: 0,
	}

	for _,attribute := range message.Attributes {
		message.Header.Length += (attribute.Length + 4)
	}

	message.Padding = message.Header.Length % 4

	message.Header.Length += message.Padding

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

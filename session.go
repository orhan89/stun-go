package stun

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type SessionID uint16

func (sessionID *SessionID) Length() uint16 {
	return 2
}

func (sessionID *SessionID) Serialize() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, sessionID)
	return buffer.Bytes()
}

func (sessionID *SessionID) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("SessionID: %d\n", sessionID))
	return buffer.String()
}

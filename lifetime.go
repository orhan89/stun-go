package stun

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Lifetime uint32

func (lifetime *Lifetime) Length() uint16 {
	return 4
}

func (lifetime *Lifetime) Serialize() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, lifetime)
	return buffer.Bytes()
}

func (lifetime *Lifetime) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Lifetime: %d\n", lifetime))
	return buffer.String()
}

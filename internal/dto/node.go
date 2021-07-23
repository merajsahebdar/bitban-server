package dto

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// NodeType
type NodeType string

// String
func (nType NodeType) String() string {
	return string(nType)
}

// Node
type Node interface{}

// ToNodeIdentifier
func ToNodeIdentifier(nType NodeType, id int64) string {
	return base64.RawStdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%d", nType, id)),
	)
}

// FromNodeIdentifier
func FromNodeIdentifier(globalID string) (nType NodeType, id int64, err error) {
	var byt []byte
	byt, err = base64.RawStdEncoding.DecodeString(globalID)

	if err == nil {
		dec := strings.Split(string(byt), ":")
		nType = NodeType(dec[0])
		id, err = strconv.ParseInt(dec[1], 10, 64)
	}

	return nType, id, err
}

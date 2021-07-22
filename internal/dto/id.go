package dto

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// ToGlobalID
func ToGlobalID(resource string, id int64) string {
	return base64.RawStdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%d", resource, id)),
	)
}

// FromGlobalID
func FromGlobalID(globalID string) (resource string, id int64, err error) {
	var byt []byte
	byt, err = base64.RawStdEncoding.DecodeString(globalID)

	if err == nil {
		dec := strings.Split(string(byt), ":")
		resource = dec[0]
		id, err = strconv.ParseInt(dec[1], 10, 64)
	}

	return resource, id, err
}

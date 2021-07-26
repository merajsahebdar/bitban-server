/*
 * Copyright 2021 Meraj Sahebdar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
func FromNodeIdentifier(nIdentifier string) (nType NodeType, id int64, err error) {
	var byt []byte
	byt, err = base64.RawStdEncoding.DecodeString(nIdentifier)

	if err == nil {
		dec := strings.Split(string(byt), ":")
		nType = NodeType(dec[0])
		id, err = strconv.ParseInt(dec[1], 10, 64)
	}

	return nType, id, err
}

// MustRetrieveIdentifier
func MustRetrieveIdentifier(nIdentifier string) int64 {
	if _, id, err := FromNodeIdentifier(nIdentifier); err != nil {
		panic(err)
	} else {
		return id
	}
}

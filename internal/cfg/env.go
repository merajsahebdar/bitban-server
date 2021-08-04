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

package cfg

// EnvType Represents types for the current running environment.
type EnvType int

const (
	Test EnvType = iota + 10
	Dev
	Prod
)

// String Returns the string representation of the current environment.
func (e EnvType) String() string {
	switch e {
	case Prod:
		return "production"
	case Dev:
		return "development"
	}

	return "unknown"
}

// CurrentEnv Keeps the value of current running environment
var CurrentEnv EnvType

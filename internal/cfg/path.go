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

import (
	"fmt"
	"os"
	"strings"
)

var (
	vars = []string{
		"/var/regeet",
	}

	etcs = []string{
		"/etc/regeet",
	}
)

// statPaths
func statPaths(dir string, dirs []string, path ...string) (string, error) {
	fin := strings.Join(path, "/")

	for _, dir := range dirs {
		if _, err := os.Stat(dir + fin); err == nil {
			return dir + fin, nil
		}
	}

	return "", fmt.Errorf("not found %s: %s", dir, fin)
}

// GetVarPath
func GetVarPath(path ...string) (string, error) {
	if found, err := statPaths("var", vars, path...); err != nil {
		return "", err
	} else {
		return found, nil
	}
}

// GetEtcPath
func GetEtcPath(path ...string) (string, error) {
	if found, err := statPaths("etc", etcs, path...); err != nil {
		return "", err
	} else {
		return found, nil
	}
}

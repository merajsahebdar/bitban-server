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

package exec

import (
	"io"
	goexec "os/exec"
)

// Create
func Create(bin string, args ...string) (cmd *goexec.Cmd, stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	cmd = goexec.Command(bin, args...)

	if stdin, err = cmd.StdinPipe(); err != nil {
		return nil, nil, nil, nil, err
	}

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, nil, nil, nil, err
	}

	if stderr, err = cmd.StderrPipe(); err != nil {
		return nil, nil, nil, nil, err
	}

	return cmd, stdin, stdout, stderr, nil
}

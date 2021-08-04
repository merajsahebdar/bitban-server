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

package facade

import (
	"context"
	"testing"
)

func TestRepo(t *testing.T) {
	t.Run("repository", func(t *testing.T) {
		testCtx := context.Background()
		testRepo := "justice"

		t.Run("create", func(t *testing.T) {
			if _, err := CreateRepoByName(testCtx, testRepo); err != nil {
				t.Errorf("failed to create the repository: %s", err.Error())
			}
		})

		t.Run("read", func(t *testing.T) {
			if _, err := GetRepoByName(testCtx, testRepo); err != nil {
				t.Errorf("got an unexpected error: %s", err.Error())
			}
		})
	})
}

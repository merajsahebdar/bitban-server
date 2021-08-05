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

package entity

import (
	"github.com/uptrace/bun"
	"github.com/volatiletech/null/v8"
)

// Policy
type Policy struct {
	bun.BaseModel `bun:"policies,select:policies,alias:policy"`
	ID            int64       `bun:"id"`
	Ptype         string      `bun:"ptype"`
	V0            string      `bun:"v0"`
	V1            string      `bun:"v1"`
	V2            string      `bun:"v2"`
	V3            null.String `bun:"v3"`
	V4            null.String `bun:"v4"`
	V5            null.String `bun:"v5"`
}

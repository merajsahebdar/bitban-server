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
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/volatiletech/null/v8"
)

// User
type User struct {
	bun.BaseModel `bun:"users,select:users,alias:user"`
	DomainID      int64       `bun:"domain_id,pk"`
	DomainType    string      `bun:"domain_type"`
	CreatedAt     time.Time   `bun:"created_at"`
	UpdatedAt     time.Time   `bun:"updated_at"`
	RemovedAt     null.Time   `bun:"removed_at"`
	Password      null.String `bun:"password"`
	IsActive      bool        `bun:"is_active"`
	IsBanned      bool        `bun:"is_banned"`
	Domain        *Domain     `bun:"rel:belongs-to,join:domain_id=id"`
	Emails        []*Email    `bun:"rel:has-many,join:domain_id=user_id"`
	Tokens        []*Token    `bun:"rel:has-many,join:domain_id=user_id"`
}

// String
func (e *User) String() string {
	return fmt.Sprintf("User<%d, %t, %t>", e.DomainID, e.IsActive, e.IsBanned)
}

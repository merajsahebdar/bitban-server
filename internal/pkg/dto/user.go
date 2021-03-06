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
	"time"

	"github.com/volatiletech/null/v8"
	"bitban.io/server/internal/pkg/orm/entity"
)

// UserNodeType
const UserNodeType NodeType = "User"

// User
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RemovedAt null.Time `json:"removedAt"`
	IsActive  bool      `json:"isActive"`
	IsBanned  bool      `json:"isBanned"`
}

// IsNode
func (User) IsNode() {}

// UserFrom Returns an instance of dto: `User` from its entity.
func UserFrom(user *entity.User) *User {
	if user != nil {
		return &User{
			ID:        ToNodeIdentifier(UserNodeType, user.DomainID),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			RemovedAt: user.RemovedAt,
			IsActive:  user.IsActive,
			IsBanned:  user.IsBanned,
		}
	}

	return nil
}

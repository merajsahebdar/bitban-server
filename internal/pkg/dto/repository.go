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
	"regeet.io/api/internal/pkg/orm/entity"
)

// RepositoryNodeType
const RepositoryNodeType NodeType = "Repository"

// Repository
type Repository struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RemovedAt null.Time `json:"removedAt"`
	Address   string    `json:"address"`
}

// IsNode
func (Repository) IsNode() {}

// RepositoryFrom Returns an instance of dto: `Repository` from its entity.
func RepositoryFrom(repository *entity.Repository) *Repository {
	if repository != nil {
		return &Repository{
			ID:        ToNodeIdentifier(RepositoryNodeType, repository.ID),
			CreatedAt: repository.CreatedAt,
			UpdatedAt: repository.UpdatedAt,
			RemovedAt: repository.RemovedAt,
			Address:   repository.Address,
		}
	}

	return nil
}
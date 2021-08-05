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
	"time"

	"github.com/uptrace/bun"
	"github.com/volatiletech/null/v8"
)

// Repository
type Repository struct {
	bun.BaseModel `bun:"repositories,select:repositories,alias:repository"`
	ID            int64      `bun:"id"`
	CreatedAt     time.Time  `bun:"created_at"`
	UpdatedAt     time.Time  `bun:"updated_at"`
	RemovedAt     null.Time  `bun:"removed_at"`
	Address       string     `bun:"address"`
	DomainID      null.Int64 `bun:"domain_id"`
	Domain        *Domain    `bun:"rel:belongs-to,join:domain_id=id"`
}

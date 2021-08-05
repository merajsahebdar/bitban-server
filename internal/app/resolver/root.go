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

package resolver

import (
	"github.com/go-playground/validator/v10"
	"regeet.io/api/internal/app/controller"
	"regeet.io/api/internal/pkg/schema"
)

type (
	// rootResolver
	rootResolver struct {
		validate          *validator.Validate
		accountController *controller.Account
	}

	// queryResolver
	queryResolver struct {
		*rootResolver
	}

	// mutationResolver
	mutationResolver struct {
		*rootResolver
	}
)

// Query
func (r *rootResolver) Query() schema.QueryResolver {
	return &queryResolver{
		rootResolver: r,
	}
}

// Mutation
func (r *rootResolver) Mutation() schema.MutationResolver {
	return &mutationResolver{
		rootResolver: r,
	}
}

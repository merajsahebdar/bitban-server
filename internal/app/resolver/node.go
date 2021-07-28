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
	"context"

	"regeet.io/api/internal/pkg/dto"
	"regeet.io/api/internal/pkg/fault"
)

// Node
func (r *queryResolver) Node(ctx context.Context, nIdentifier string) (node dto.Node, err error) {
	var id int64
	var nType dto.NodeType

	if nType, id, err = dto.FromNodeIdentifier(nIdentifier); err != nil {
		return nil, NotFoundErrorFrom(err)
	}

	switch nType {
	case dto.UserNodeType:
		node, err = r.accountController.GetUser(ctx, id)
	}

	if err == nil {
		return node, nil
	}

	switch {
	case fault.IsUnauthenticatedError(err):
		return nil, AuthenticationErrorFrom(err)
	case fault.IsForbiddenError(err):
		return nil, ForbiddenErrorFrom(err)
	case fault.IsResourceNotFoundError(err):
		return nil, NotFoundErrorFrom(err)
	default:
		panic(err)
	}
}

package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/fault"
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

package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/dto"
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
	case common.IsUnauthenticatedError(err):
		return nil, AuthenticationErrorFrom(err)
	case common.IsForbiddenError(err):
		return nil, ForbiddenErrorFrom(err)
	case common.IsResourceNotFoundError(err):
		return nil, NotFoundErrorFrom(err)
	default:
		panic(err)
	}
}

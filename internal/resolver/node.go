package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/facade"
)

// Node
func (*queryResolver) Node(ctx context.Context, nIdentifier string) (dto.Node, error) {
	var err error
	var id int64
	var nType dto.NodeType

	if nType, id, err = dto.FromNodeIdentifier(nIdentifier); err != nil {
		return nil, NotFoundError(nil)
	}

	switch nType {
	case dto.UserNodeType:
		if account, err := facade.GetAccountByUserID(ctx, id); common.IsResourceNotFoundError(err) {
			return nil, nil
		} else {
			return dto.UserFrom(
				account.GetUser(),
			), nil
		}
	}

	panic(err)
}

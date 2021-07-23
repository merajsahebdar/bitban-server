package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/dto"
)

// User Returns an existing user using its identifier.
func (*queryResolver) User(context.Context, dto.UserFilter) (*dto.User, error) {
	panic("not implemented")
}

// Profile
func (*userResolver) Profile(context.Context, *dto.User) (*dto.UserProfile, error) {
	panic("not implemeted")
}

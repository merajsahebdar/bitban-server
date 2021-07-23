package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/dto"
)

// Profile
func (*userResolver) Profile(context.Context, *dto.User) (*dto.UserProfile, error) {
	panic("not implemeted")
}

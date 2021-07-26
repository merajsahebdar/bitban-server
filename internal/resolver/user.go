package resolver

import (
	"context"

	"regeet.io/api/internal/dto"
)

// Profile
func (*userResolver) Profile(context.Context, *dto.User) (*dto.UserProfile, error) {
	panic("not implemeted")
}

package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/facade"
)

// User Returns an existing user using its identifier.
func (*queryResolver) User(context.Context, dto.UserFilter) (*dto.User, error) {
	panic("not implemented")
}

// SignIn
func (r *mutationResolver) SignIn(ctx context.Context, input dto.SignInInput) (*dto.Auth, error) {
	if err := r.validate.Struct(input); err != nil {
		return nil, UserInputErrorFrom(
			common.UserInputErrorFrom(err),
		)
	}

	if account, err := facade.GetAccountByPassword(
		ctx,
		input,
	); common.IsNonUserInputError(err) {
		panic(err)
	} else if common.IsUserInputError(err) {
		return nil, UserInputErrorFrom(err)
	} else {
		var accessToken string
		if accessToken, err = account.CreateAccessToken(); err != nil {
			panic(err)
		}

		return &dto.Auth{
			AccessToken: accessToken,
			User:        dto.UserFrom(account.GetUser()),
		}, nil
	}
}

// SignUp
func (r *mutationResolver) SignUp(ctx context.Context, input dto.SignUpInput) (*dto.Auth, error) {
	if err := r.validate.Struct(input); err != nil {
		return nil, UserInputErrorFrom(
			common.UserInputErrorFrom(err),
		)
	}

	if account, err := facade.CreateAccount(ctx, input); err != nil {
		panic(err)
	} else {
		var accessToken string
		if accessToken, err = account.CreateAccessToken(); err != nil {
			panic(err)
		}

		return &dto.Auth{
			AccessToken: accessToken,
			User:        dto.UserFrom(account.GetUser()),
		}, nil
	}
}

// RefreshToken
func (*mutationResolver) RefreshToken(context.Context) (string, error) {
	panic("not implemented")
}

// Profile
func (*userResolver) Profile(context.Context, *dto.User) (*dto.UserProfile, error) {
	panic("not implemeted")
}

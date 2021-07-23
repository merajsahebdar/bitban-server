package resolver

import (
	"context"

	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/facade"
)

// createAccessTokenByAccount
func createAccessTokenByAccoount(account *facade.Account) string {
	if accessToken, err := account.CreateAccessToken(); err != nil {
		panic(err)
	} else {
		return accessToken
	}
}

// createAuthTokensByAccount
func createAuthTokensByAccount(account *facade.Account) (string, string) {
	var err error

	var refreshToken string
	if refreshToken, err = account.CreateRefreshToken(); err != nil {
		panic(err)
	}

	return refreshToken, createAccessTokenByAccoount(account)
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
		_, accessToken := createAuthTokensByAccount(account)

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
		_, accessToken := createAuthTokensByAccount(account)

		return &dto.Auth{
			AccessToken: accessToken,
			User:        dto.UserFrom(account.GetUser()),
		}, nil
	}
}

// RefreshToken
func (*mutationResolver) RefreshToken(ctx context.Context) (string, error) {
	if account, err := facade.GetAccountByUser(
		ctx,
		common.GetContextAuthorizedUser(ctx),
	); err != nil {
		panic(err)
	} else {
		return createAccessTokenByAccoount(
			account,
		), nil
	}
}

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

	"regeet.io/api/internal/pkg/auth"
	"regeet.io/api/internal/pkg/dto"
	"regeet.io/api/internal/pkg/facade"
	"regeet.io/api/internal/pkg/fault"
)

//
// TODO: delegate logic to the controller layer
//

// createAccessTokenByAccount
func createAccessTokenByAccoount(account *facade.Account) string {
	if accessToken, err := account.CreateAccessToken(); err != nil {
		panic(err)
	} else {
		return accessToken
	}
}

// createAuthTokensByAccount
func createAuthTokensByAccount(ctx context.Context, account *facade.Account) (string, string) {
	var err error

	var refreshToken string
	if refreshToken, err = account.CreateRefreshToken(); err != nil {
		panic(err)
	}

	auth.SetRefreshTokenCookie(ctx, refreshToken)

	return refreshToken, createAccessTokenByAccoount(account)
}

// SignIn
func (r *mutationResolver) SignIn(ctx context.Context, input dto.SignInInput) (*dto.Auth, error) {
	if err := r.validate.Struct(input); err != nil {
		return nil, UserInputErrorFrom(
			fault.UserInputErrorFrom(err),
		)
	}

	if account, err := facade.GetAccountByPassword(
		ctx,
		input,
	); fault.IsNonUserInputError(err) {
		panic(err)
	} else if fault.IsUserInputError(err) {
		return nil, UserInputErrorFrom(err)
	} else {
		_, accessToken := createAuthTokensByAccount(ctx, account)

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
			fault.UserInputErrorFrom(err),
		)
	}

	if account, err := facade.CreateAccount(ctx, input); err != nil {
		panic(err)
	} else {
		_, accessToken := createAuthTokensByAccount(ctx, account)

		return &dto.Auth{
			AccessToken: accessToken,
			User:        dto.UserFrom(account.GetUser()),
		}, nil
	}
}

// RefreshToken
func (*mutationResolver) RefreshToken(ctx context.Context) (string, error) {
	if account, err := facade.GetAccountByRefreshToken(
		ctx,
	); err != nil {
		return "", AuthenticationErrorFrom(err)
	} else {
		return createAccessTokenByAccoount(
			account,
		), nil
	}
}

// CreateRepository
func (r *mutationResolver) CreateRepository(ctx context.Context, input dto.CreateRepositoryInput) (*dto.Repository, error) {
	if repository, err := r.
		repoController.
		CreateRepository(ctx, input); err != nil {
		switch {
		case fault.IsUnauthenticatedError(err):
			return nil, AuthenticationErrorFrom(err)
		case fault.IsUserInputError(err):
			return nil, UserInputErrorFrom(err)
		default:
			panic(err)
		}
	} else {
		return dto.RepositoryFrom(repository), nil
	}
}

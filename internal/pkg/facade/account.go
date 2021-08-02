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

package facade

import (
	"context"
	"fmt"
	"time"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/uptrace/bun"
	"github.com/volatiletech/null/v8"
	"regeet.io/api/internal/cfg"
	"regeet.io/api/internal/pkg/auth"
	"regeet.io/api/internal/pkg/db"
	"regeet.io/api/internal/pkg/db/orm"
	"regeet.io/api/internal/pkg/dto"
	"regeet.io/api/internal/pkg/fault"
	"regeet.io/api/internal/pkg/jwt"
	"regeet.io/api/internal/pkg/util"
)

//
// TODO: move token signing into auth package
//

// defaultDomain
const defaultDomain = "_"

type (
	// Account
	Account struct {
		ctx  context.Context
		user *orm.User
	}
)

// GetUser
func (f *Account) GetUser() *orm.User {
	return f.user
}

// CheckPermission
//
// Errors:
//   - common.ErrForbidden if the authorized user does not have access to the resource
func (f *Account) CheckPermission(obj string, act string) error {
	if ok, err := auth.GetEnforcerInstance().Enforce(
		fmt.Sprintf("/users/%d", f.user.ID),
		defaultDomain,
		obj,
		act,
	); err != nil || !ok {
		return fault.ErrForbidden
	}

	return nil
}

// CreateAccessToken
func (f *Account) CreateAccessToken() (accessToken string, err error) {
	comp := jwt.GetJwtInstance()

	currTime := time.Now().In(time.UTC)
	claims := &gojwt.StandardClaims{
		Subject:  dto.ToNodeIdentifier(dto.UserNodeType, f.user.ID),
		IssuedAt: currTime.Unix(),
		ExpiresAt: currTime.Add(
			time.Duration(cfg.Cog.Security.AccessTokenExpiresAt) * time.Minute,
		).Unix(),
	}

	accessToken, err = comp.SignToken(claims)
	return accessToken, err
}

// CreateRefreshToken
func (f *Account) CreateRefreshToken() (refreshToken string, err error) {
	comp := jwt.GetJwtInstance()

	userToken := &orm.UserToken{
		Meta:   struct{}{},
		UserID: null.Int64From(f.user.ID),
	}
	if _, err = db.GetBunInstance().
		NewInsert().
		Model(userToken).
		Column("meta", "user_id").
		Exec(f.ctx); err != nil {
		return "", err
	}

	currTime := time.Now().In(time.UTC)
	expiresAt := currTime.Add(
		time.Duration(cfg.Cog.Security.RefreshTokenExpiresAt) * time.Minute,
	)
	claims := &gojwt.StandardClaims{
		Id:        dto.ToNodeIdentifier(dto.UserTokenNodeType, userToken.ID),
		Subject:   dto.ToNodeIdentifier(dto.UserNodeType, f.user.ID),
		IssuedAt:  currTime.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	if refreshToken, err = comp.SignToken(claims); err != nil {
		return "", err
	}

	auth.SetRefreshTokenCookie(f.ctx, refreshToken)

	return refreshToken, nil
}

// GetAccountByPassword
//
// If was not able to find the corresponding account, returns `fault.ErrUserInput`.
func GetAccountByPassword(ctx context.Context, input dto.SignInInput) (*Account, error) {
	userPrimaryEmail := new(orm.UserEmail)
	if err := db.GetBunInstance().
		NewSelect().
		Model(userPrimaryEmail).
		Relation("User", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("? = ?", bun.Ident("user.is_active"), true).
				Where("? = ?", bun.Ident("user.is_banned"), false).
				Where("? IS NULL", bun.Ident("user.removed_at"), nil)
		}).
		Relation("User.Profile", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("? IS NULL", bun.Ident("user__profile.removed_at"))
		}).
		Where("? = ?", bun.Ident("user_email.address"), input.Identifier).
		Where("? = ?", bun.Ident("user_email.is_primary"), true).
		Where("? = ?", bun.Ident("user_email.is_verified"), true).
		Where("? IS NULL", bun.Ident("user_email.removed_at")).
		Limit(1).
		Scan(ctx); fault.IsNonResourceNotFoundError(err) {
		return nil, err
	} else if fault.IsResourceNotFoundError(err) {
		return nil, fault.ErrUserInput
	}

	user := userPrimaryEmail.User
	if user == nil || user.Password.IsZero() || !util.ComparePassword(user.Password.String, input.Password) {
		return nil, fault.ErrUserInput
	}

	return &Account{
		ctx:  ctx,
		user: user,
	}, nil
}

// GetAccountByUser
func GetAccountByUser(ctx context.Context, user *orm.User) (account *Account, err error) {
	account = &Account{
		ctx:  ctx,
		user: user,
	}
	return account, err
}

// GetAccountByUserId
func GetAccountByUserId(ctx context.Context, id int64) (*Account, error) {
	user := new(orm.User)
	if err := db.GetBunInstance().
		NewSelect().
		Model(user).
		Where("? = ?", bun.Ident("user.id"), id).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	} else {
		return &Account{
			ctx:  ctx,
			user: user,
		}, nil
	}
}

// CreateAccount
func CreateAccount(ctx context.Context, input dto.SignUpInput) (account *Account, err error) {
	var tx bun.Tx
	if tx, err = db.GetBunInstance().BeginTx(ctx, nil); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && account == nil {
			tx.Rollback()
		}
	}()

	//
	// Create User

	var hashedPassword string
	if hashedPassword, err = util.HashPassword(input.Password); err != nil {
		return nil, err
	}

	user := &orm.User{
		Password: null.StringFrom(hashedPassword),
		IsActive: true,
		IsBanned: false,
	}
	if _, err = tx.NewInsert().
		Model(user).
		Column("password", "is_active", "is_banned").
		Returning("id", "created_at", "updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}

	//
	// Create User's Email

	userEmail := &orm.UserEmail{
		Address:    input.PrimaryEmail.Address,
		IsVerified: true,
		IsPrimary:  true,
		UserID:     null.Int64From(user.ID),
	}
	if _, err = tx.NewInsert().
		Model(userEmail).
		Column("address", "is_verified", "is_primary", "user_id").
		Exec(ctx); err != nil {
		return nil, err
	}

	//
	// Create User's Profile

	userProfile := &orm.UserProfile{
		Name:   input.Profile.Name,
		Meta:   struct{}{},
		UserID: null.Int64From(user.ID),
	}
	if _, err = tx.NewInsert().
		Model(userProfile).
		Column("name", "meta", "user_id").
		Exec(ctx); err != nil {
		return nil, err
	}

	user.Emails = append(user.Emails, userEmail)
	user.Profile = userProfile

	//
	// Last Step!

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Grant Permissions
	sub := fmt.Sprintf("/users/%d", user.ID)
	if _, err = auth.GetEnforcerInstance().AddNamedPolicies(
		"p",
		[][]string{
			{sub, defaultDomain, fmt.Sprintf("/users/%d", user.ID), ".*"},
			{sub, defaultDomain, fmt.Sprintf("/users/%d/*", user.ID), ".*"},
		},
	); err != nil {
		return nil, err
	}

	account = &Account{
		ctx:  ctx,
		user: user,
	}

	return account, nil
}

// GetAccountByAccessToken
//
// ErrorsRef:
//   - auth.GetContextAccessTokenClaims
//   - facade.GetAccountByUserId
func GetAccountByAccessToken(ctx context.Context) (*Account, error) {
	if claims, err := auth.GetContextAccessTokenClaims(ctx); err != nil {
		return nil, err
	} else {
		return GetAccountByUserId(
			ctx,
			dto.MustRetrieveIdentifier(claims.Subject),
		)
	}
}

// GetAccountByRefreshToken
//
// ErrorsRef:
//   - auth.GetContextRefreshTokenClaims
//   - facade.GetAccountByUserId
func GetAccountByRefreshToken(ctx context.Context) (*Account, error) {
	if claims, err := auth.GetContextRefreshTokenClaims(ctx); err != nil {
		return nil, err
	} else {
		return GetAccountByUserId(
			ctx,
			dto.MustRetrieveIdentifier(claims.Subject),
		)
	}
}

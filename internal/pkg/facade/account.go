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
	"regeet.io/api/internal/pkg/dto"
	"regeet.io/api/internal/pkg/fault"
	"regeet.io/api/internal/pkg/jwt"
	"regeet.io/api/internal/pkg/orm"
	"regeet.io/api/internal/pkg/orm/entity"
	"regeet.io/api/internal/pkg/util"
)

//
// TODO: move token signing into auth package
//

// appDomain
const appDomain = "_"

type (
	// Account
	Account struct {
		ctx  context.Context
		user *entity.User
	}
)

// GetUser
func (f *Account) GetUser() *entity.User {
	return f.user
}

// GetDomain
func (f *Account) GetDomain() *entity.Domain {
	return f.user.Domain
}

// CheckPermission
//
// Errors:
//   - fault.ErrForbidden if the authorized user does not have access to the resource
func (f *Account) CheckPermission(obj string, act string) error {
	return f.CheckPermissionIn(appDomain, obj, act)
}

// CheckPermissionIn
func (f *Account) CheckPermissionIn(dom string, obj string, act string) error {
	if ok, err := auth.
		GetEnforcerInstance().
		Enforce(
			fmt.Sprintf("/users/%d", f.user.DomainID),
			dom,
			obj,
			act,
		); err != nil || !ok {
		return fault.ErrForbidden
	}

	return nil
}

// CreateAccessToken
func (f *Account) CreateAccessToken() (accessToken string, err error) {
	currTime := time.Now().In(time.UTC)
	claims := &gojwt.StandardClaims{
		Subject:  dto.ToNodeIdentifier(dto.UserNodeType, f.user.DomainID),
		IssuedAt: currTime.Unix(),
		ExpiresAt: currTime.Add(
			time.Duration(cfg.Cog.Security.AccessTokenExpiresAt) * time.Minute,
		).Unix(),
	}

	accessToken, err = jwt.
		GetJwtInstance().
		SignToken(claims)
	return accessToken, err
}

// CreateRefreshToken
func (f *Account) CreateRefreshToken() (refreshToken string, err error) {
	token := &entity.Token{
		Meta:   struct{}{},
		UserID: null.Int64From(f.user.DomainID),
	}
	if _, err = orm.GetBunInstance().
		NewInsert().
		Model(token).
		Column("meta", "user_id").
		Exec(f.ctx); err != nil {
		return "", err
	}

	currTime := time.Now().In(time.UTC)
	expiresAt := currTime.Add(
		time.Duration(cfg.Cog.Security.RefreshTokenExpiresAt) * time.Minute,
	)
	claims := &gojwt.StandardClaims{
		Id:        dto.ToNodeIdentifier(dto.TokenNodeType, token.ID),
		Subject:   dto.ToNodeIdentifier(dto.UserNodeType, f.user.DomainID),
		IssuedAt:  currTime.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	if refreshToken, err = jwt.
		GetJwtInstance().
		SignToken(claims); err != nil {
		return "", err
	}

	return refreshToken, nil
}

// GetAccountByPassword
//
// Errors:
//   - fault.ErrUserInput if was not able to find the corresponding account
func GetAccountByPassword(ctx context.Context, input dto.SignInInput) (*Account, error) {
	primaryEmail := new(entity.Email)
	if err := orm.GetBunInstance().
		NewSelect().
		Model(primaryEmail).
		Relation("User", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("? = ?", bun.Ident("user.is_active"), true).
				Where("? = ?", bun.Ident("user.is_banned"), false).
				Where("? IS NULL", bun.Ident("user.removed_at"), nil)
		}).
		Relation("User.Domain", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("? IS NULL", bun.Ident("user__domain.removed_at"))
		}).
		Where("? = ?", bun.Ident("email.address"), input.Identifier).
		Where("? = ?", bun.Ident("email.is_primary"), true).
		Where("? = ?", bun.Ident("email.is_verified"), true).
		Where("? IS NULL", bun.Ident("email.removed_at")).
		Limit(1).
		Scan(ctx); fault.IsNonResourceNotFoundError(err) {
		return nil, err
	} else if fault.IsResourceNotFoundError(err) {
		return nil, fault.ErrUserInput
	}

	user := primaryEmail.User
	if user == nil || user.Password.IsZero() || !util.ComparePassword(user.Password.String, input.Password) {
		return nil, fault.ErrUserInput
	}

	return &Account{
		ctx:  ctx,
		user: user,
	}, nil
}

// GetAccountByUserId
func GetAccountByUserId(ctx context.Context, id int64) (*Account, error) {
	user := new(entity.User)
	if err := orm.GetBunInstance().
		NewSelect().
		Model(user).
		Relation("Domain").
		Where("? = ?", bun.Ident("user.domain_id"), id).
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
	if tx, err = orm.GetBunInstance().BeginTx(ctx, nil); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && account == nil {
			tx.Rollback()
		}
	}()

	//
	// Create Domain

	domain := &entity.Domain{
		Type:    "user",
		Name:    input.Domain.Name,
		Address: input.Domain.Address,
		Meta:    struct{}{},
	}
	if _, err := tx.NewInsert().
		Model(domain).
		Column("type", "name", "address", "meta").
		Returning("id", "created_at", "updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}

	//
	// Create User

	var hashedPassword string
	if hashedPassword, err = util.HashPassword(input.Password); err != nil {
		return nil, err
	}

	user := &entity.User{
		DomainID:   domain.ID,
		DomainType: domain.Type,
		Password:   null.StringFrom(hashedPassword),
		IsActive:   true,
		IsBanned:   false,
	}
	if _, err = tx.NewInsert().
		Model(user).
		Column("domain_id", "domain_type", "password", "is_active", "is_banned").
		Returning("created_at", "updated_at").
		Exec(ctx); err != nil {
		return nil, err
	}

	//
	// Create User's Email

	email := &entity.Email{
		Address:    input.PrimaryEmail.Address,
		IsVerified: true,
		IsPrimary:  true,
		UserID:     null.Int64From(user.DomainID),
	}
	if _, err = tx.NewInsert().
		Model(email).
		Column("address", "is_verified", "is_primary", "user_id").
		Exec(ctx); err != nil {
		return nil, err
	}

	user.Emails = append(user.Emails, email)

	//
	// Last Step!

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Grant Permissions
	sub := fmt.Sprintf("/users/%d", user.DomainID)
	if _, err = auth.GetEnforcerInstance().AddNamedPolicies(
		"p",
		[][]string{
			{sub, appDomain, fmt.Sprintf("/users/%d", user.DomainID), ".*"},
			{sub, appDomain, fmt.Sprintf("/users/%d/*", user.DomainID), ".*"},
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

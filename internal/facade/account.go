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
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"regeet.io/api/internal/auth"
	"regeet.io/api/internal/component"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/db"
	"regeet.io/api/internal/dto"
	"regeet.io/api/internal/fault"
	"regeet.io/api/internal/orm"
	"regeet.io/api/internal/util"
)

//
// TODO: move token signing into auth package
//

// defaultDomain
const defaultDomain = "_"

type (
	// Account
	Account struct {
		ctx         context.Context
		user        *orm.User
		userEmail   *orm.UserEmail
		userProfile *orm.UserProfile
	}

	// accountBinder
	accountBinder struct {
		User        orm.User        `boil:"users,bind"`
		UserEmail   orm.UserEmail   `boil:"user_emails,bind"`
		UserProfile orm.UserProfile `boil:"user_profiles,bind"`
	}
)

var (
	// accountColumnsSelection
	accountColumnsSelection = qm.Select(
		orm.UserTableColumns.ID,
		orm.UserTableColumns.Password,
		orm.UserEmailTableColumns.ID,
		orm.UserEmailTableColumns.Address,
		orm.UserEmailTableColumns.UserID,
		orm.UserProfileTableColumns.ID,
		orm.UserProfileTableColumns.Name,
		orm.UserProfileTableColumns.Meta,
		orm.UserProfileTableColumns.UserID,
	)
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
	comp := component.GetJwtInstance()

	currTime := time.Now().In(time.UTC)
	claims := &jwt.StandardClaims{
		Subject:  dto.ToNodeIdentifier(dto.UserNodeType, f.user.ID),
		IssuedAt: currTime.Unix(),
		ExpiresAt: currTime.Add(
			time.Duration(conf.Cog.Security.AccessTokenExpiresAt) * time.Minute,
		).Unix(),
	}

	accessToken, err = comp.SignToken(claims)
	return accessToken, err
}

// CreateRefreshToken
func (f *Account) CreateRefreshToken() (refreshToken string, err error) {
	comp := component.GetJwtInstance()

	userToken := &orm.UserToken{
		Meta:   []byte(`{}`),
		UserID: null.Int64From(f.user.ID),
	}
	if err = userToken.Insert(f.ctx, db.GetDbInstance(), boil.Infer()); err != nil {
		return "", err
	}

	currTime := time.Now().In(time.UTC)
	expiresAt := currTime.Add(
		time.Duration(conf.Cog.Security.RefreshTokenExpiresAt) * time.Minute,
	)
	claims := &jwt.StandardClaims{
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
// If was not able to find the corresponding account, returns `common.ErrUserInput`.
func GetAccountByPassword(ctx context.Context, input dto.SignInInput) (*Account, error) {
	var err error

	var binder accountBinder
	if err = orm.NewQuery(
		accountColumnsSelection,
		qm.From(`"users"`),
		qm.InnerJoin(`"user_emails" ON "user_emails"."user_id" = "users"."id"`),
		qm.InnerJoin(`"user_profiles" ON "user_profiles"."user_id" = "users"."id"`),
		orm.UserWhere.IsActive.EQ(true),
		orm.UserWhere.IsBanned.EQ(false),
		orm.UserWhere.RemovedAt.IsNull(),
		orm.UserEmailWhere.Address.EQ(input.Identifier),
		orm.UserEmailWhere.IsPrimary.EQ(true),
		orm.UserEmailWhere.IsVerified.EQ(true),
		orm.UserEmailWhere.RemovedAt.IsNull(),
	).Bind(ctx, db.GetDbInstance(), &binder); err != nil {
		return nil, err
	} else if binder.User == (orm.User{}) {
		return nil, fault.ErrUserInput
	}

	if binder.User.Password.IsZero() || !util.ComparePassword(binder.User.Password.String, input.Password) {
		return nil, fault.ErrUserInput
	}

	return &Account{
		ctx:         ctx,
		user:        &binder.User,
		userEmail:   &binder.UserEmail,
		userProfile: &binder.UserProfile,
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
	if user, err := orm.FindUser(
		ctx,
		db.GetDbInstance(),
		id,
	); err != nil {
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
	var tx *sql.Tx
	if tx, err = db.GetDbInstance().BeginTx(ctx, nil); err != nil {
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
	if err = user.Insert(ctx, tx, boil.Infer()); err != nil {
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
	if err = userEmail.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}

	//
	// Create User's Profile

	userProfile := &orm.UserProfile{
		Name:   input.Profile.Name,
		Meta:   []byte(`{}`),
		UserID: null.Int64From(user.ID),
	}
	if err = userProfile.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}

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
		ctx:         ctx,
		user:        user,
		userEmail:   userEmail,
		userProfile: userProfile,
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

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
	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/conf"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/orm"
)

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
	if ok, err := common.GetEnforcerInstance().Enforce(
		fmt.Sprintf("/users/%d", f.user.ID),
		common.DefaultUserDomain,
		obj,
		act,
	); err != nil || !ok {
		return common.ErrForbidden
	}

	return nil
}

// CreateAccessToken
func (f *Account) CreateAccessToken() (accessToken string, err error) {
	comp := common.GetJwtInstance()

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
	db := common.GetContextDb(f.ctx)
	comp := common.GetJwtInstance()

	userToken := &orm.UserToken{
		Meta:   []byte(`{}`),
		UserID: null.Int64From(f.user.ID),
	}
	if err = userToken.Insert(f.ctx, db, boil.Infer()); err != nil {
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

	common.SetRefreshTokenCookie(f.ctx, refreshToken)

	return refreshToken, nil
}

// GetAccountByPassword
//
// If was not able to find the corresponding account, returns `common.ErrUserInput`.
func GetAccountByPassword(ctx context.Context, input dto.SignInInput) (*Account, error) {
	db := common.GetContextDb(ctx)

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
	).Bind(ctx, db, &binder); err != nil {
		return nil, err
	} else if binder.User == (orm.User{}) {
		return nil, common.ErrUserInput
	}

	if binder.User.Password.IsZero() || !common.ComparePassword(binder.User.Password.String, input.Password) {
		return nil, common.ErrUserInput
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
		common.GetContextDb(ctx),
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
	db := common.GetContextDb(ctx)

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
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
	if hashedPassword, err = common.HashPassword(input.Password); err != nil {
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
	dom := common.DefaultUserDomain
	if _, err = common.GetEnforcerInstance().AddNamedPolicies(
		"p",
		[][]string{
			{sub, dom, fmt.Sprintf("/users/%d", user.ID), ".*"},
			{sub, dom, fmt.Sprintf("/users/%d/*", user.ID), ".*"},
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
//   - common.GetContextAccessTokenClaims
//   - facade.GetAccountByUserId
func GetAccountByAccessToken(ctx context.Context) (*Account, error) {
	if claims, err := common.GetContextAccessTokenClaims(ctx); err != nil {
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
//   - common.GetContextRefreshTokenClaims
//   - facade.GetAccountByUserId
func GetAccountByRefreshToken(ctx context.Context) (*Account, error) {
	if claims, err := common.GetContextRefreshTokenClaims(ctx); err != nil {
		return nil, err
	} else {
		return GetAccountByUserId(
			ctx,
			dto.MustRetrieveIdentifier(claims.Subject),
		)
	}
}

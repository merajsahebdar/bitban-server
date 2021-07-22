package facade

import (
	"context"
	"database/sql"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/orm"
)

type (
	// Account
	Account struct {
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

	// defaultUserProfileMeta
	defaultUserProfileMeta = []byte(`{}`)
)

// GetUser
func (f *Account) GetUser() *orm.User {
	return f.user
}

// GetAccountByPassword
//
// If was not able to find the corresponding account, returns `common.ErrUserInput`.
func GetAccountByPassword(ctx context.Context, input dto.SignInInput) (*Account, error) {
	db := common.GetContextDB(ctx)

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
		user:        &binder.User,
		userEmail:   &binder.UserEmail,
		userProfile: &binder.UserProfile,
	}, nil
}

// CreateAccount
func CreateAccount(ctx context.Context, input dto.SignUpInput) (account *Account, err error) {
	db := common.GetContextDB(ctx)

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
		Meta:   defaultUserProfileMeta,
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

	account = &Account{
		user:        user,
		userEmail:   userEmail,
		userProfile: userProfile,
	}

	return account, nil
}

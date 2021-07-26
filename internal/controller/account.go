package controller

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"go.uber.org/fx"
	"regeet.io/api/internal/auth"
	"regeet.io/api/internal/dto"
	"regeet.io/api/internal/facade"
	"regeet.io/api/internal/fault"
)

// Account
type Account struct {
	enforcer *casbin.Enforcer
}

// GetUser
//
// Errors:
//   - fault.Unauthenticated, if the request is not authorized
//   - fault.ErrForbidden, if the authorized user does not have access to the resource
func (c *Account) GetUser(ctx context.Context, id int64) (*dto.User, error) {
	//
	// Check Permission

	if currAccount, err := facade.GetAccountByAccessToken(ctx); err != nil {
		return nil, fault.ErrUnauthenticated
	} else {
		if err := currAccount.CheckPermission(
			fmt.Sprintf("/users/%d", id),
			"read",
		); err != nil {
			return nil, err
		}
	}

	//
	// Rerieve the User

	if account, err := facade.GetAccountByUserId(ctx, id); err != nil {
		return nil, err
	} else {
		return dto.UserFrom(
			account.GetUser(),
		), nil
	}
}

// AccountOpt
var AccountOpt = fx.Provide(newAccount)

// newAccount
func newAccount() *Account {
	return &Account{
		enforcer: auth.GetEnforcerInstance(),
	}
}

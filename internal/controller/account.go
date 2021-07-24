package controller

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"go.giteam.ir/giteam/internal/auth"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/facade"
	"go.giteam.ir/giteam/internal/fault"
	"go.uber.org/fx"
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

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

package controller

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	"go.uber.org/fx"
	"bitban.io/server/internal/pkg/auth"
	"bitban.io/server/internal/pkg/dto"
	"bitban.io/server/internal/pkg/facade"
	"bitban.io/server/internal/pkg/fault"
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

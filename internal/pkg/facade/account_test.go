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
	"testing"

	"regeet.io/api/internal/pkg/dto"
	"regeet.io/api/internal/pkg/fault"
	"syreclabs.com/go/faker"
)

func TestAccount(t *testing.T) {
	t.Run("account", func(t *testing.T) {
		ctx := context.Background()

		newInput := struct {
			password   string
			identifier string
			name       string
		}{
			password:   faker.Internet().Password(8, 10),
			identifier: faker.Internet().SafeEmail(),
			name:       faker.Name().Name(),
		}

		t.Run("sign-up-valid", func(t *testing.T) {
			if account, err := CreateAccount(ctx, dto.SignUpInput{
				Password:        newInput.password,
				PasswordConfirm: newInput.password,
				PrimaryEmail: dto.SignUpPrimaryEmailInput{
					Address: newInput.identifier,
				},
				Profile: dto.SignUpProfileInput{
					Name: newInput.name,
				},
			}); err != nil {
				t.Errorf(err.Error())
			} else {
				t.Logf("signed up as %s", account.GetUser().String())
			}
		})

		t.Run("sign-up-invalid", func(t *testing.T) {
			if _, err := CreateAccount(ctx, dto.SignUpInput{
				Password:        newInput.password,
				PasswordConfirm: newInput.password,
				PrimaryEmail: dto.SignUpPrimaryEmailInput{
					Address: newInput.identifier,
				},
				Profile: dto.SignUpProfileInput{
					Name: newInput.name,
				},
			}); fault.IsNonPqUniqueViolationError(err) {
				t.Errorf(err.Error())
			}
		})

		t.Run("sign-in-valid", func(t *testing.T) {
			if account, err := GetAccountByPassword(ctx, dto.SignInInput{
				Identifier: newInput.identifier,
				Password:   newInput.password,
			}); err != nil {
				t.Errorf(err.Error())
			} else {
				t.Logf("signed in as %s", account.GetUser().String())
			}
		})

		t.Run("sign-in-invalid-email", func(t *testing.T) {
			if _, err := GetAccountByPassword(ctx, dto.SignInInput{
				Identifier: "invalid",
				Password:   newInput.password,
			}); fault.IsNonUserInputError(err) {
				t.Errorf(err.Error())
			}
		})

		t.Run("sign-in-invalid-password", func(t *testing.T) {
			if _, err := GetAccountByPassword(ctx, dto.SignInInput{
				Identifier: newInput.identifier,
				Password:   "invalid",
			}); fault.IsNonUserInputError(err) {
				t.Errorf(err.Error())
			}
		})
	})
}

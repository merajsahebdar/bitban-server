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

package dto

// SignInInput
type SignInInput struct {
	Password   string `json:"password" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
}

// SignUpInput
type SignUpInput struct {
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`

	Domain       SignUpDomainInput       `json:"domain"`
	PrimaryEmail SignUpPrimaryEmailInput `json:"primaryEmail"`
}

// SignUpDomainInput
type SignUpDomainInput struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required,notexistsin=domains address"`
}

// SignUpPrimaryEmailInput
type SignUpPrimaryEmailInput struct {
	Address string `json:"address" validate:"required,email,notexistsin=emails address"`
}

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

// Auth
type Auth struct {
	AccessToken string `json:"accessToken"`
	User        *User  `json:"user"`
}

// SignInInput
type SignInInput struct {
	Password   string `json:"password" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
}

// SignUpInput
type SignUpInput struct {
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`

	Profile      SignUpProfileInput      `json:"profile"`
	PrimaryEmail SignUpPrimaryEmailInput `json:"primaryEmail"`
}

// SignUpPrimaryEmailInput
type SignUpPrimaryEmailInput struct {
	Address string `json:"address" validate:"required,email,notexistsin=user_emails address"`
}

// SignUpProfileInput
type SignUpProfileInput struct {
	Name string `json:"name" validate:"required"`
}

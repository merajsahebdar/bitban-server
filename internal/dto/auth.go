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

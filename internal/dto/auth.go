package dto

// Auth
type Auth struct {
	AccessToken string `json:"accessToken"`
	User        *User  `json:"user"`
}

// SignInInput
type SignInInput struct {
	Password   string `json:"password"`
	Identifier string `json:"identifier"`
}

// SignUpInput
type SignUpInput struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`

	Profile      SignUpProfileInput      `json:"profile"`
	PrimaryEmail SignUpPrimaryEmailInput `json:"primaryEmail"`
}

// SignUpPrimaryEmailInput
type SignUpPrimaryEmailInput struct {
	Address string `json:"address"`
}

// SignUpProfileInput
type SignUpProfileInput struct {
	Name string `json:"name"`
}

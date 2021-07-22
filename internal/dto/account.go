package dto

import (
	"time"

	"github.com/volatiletech/null/v8"
	"go.giteam.ir/giteam/internal/orm"
)

// User
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RemovedAt null.Time `json:"removedAt"`
	IsActive  bool      `json:"isActive"`
	IsBanned  bool      `json:"isBanned"`
}

// UserFrom Returns an instance of model: `User` from its datasource.
func UserFrom(user *orm.User) *User {
	if user != nil {
		return &User{
			ID:        ToGlobalID("User", user.ID),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			RemovedAt: user.RemovedAt,
			IsActive:  user.IsActive,
			IsBanned:  user.IsBanned,
		}
	}

	return nil
}

// UserProfile
type UserProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserProfileFrom
func UserProfileFrom(profile *orm.UserProfile) *UserProfile {
	if profile != nil {
		return &UserProfile{
			ID:   ToGlobalID("UserProfile", profile.ID),
			Name: profile.Name,
		}
	}

	return nil
}

// UserFilter
type UserFilter struct {
	ID string `json:"id" param:"id"`
}

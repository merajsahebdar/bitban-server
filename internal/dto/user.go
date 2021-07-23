package dto

import (
	"time"

	"github.com/volatiletech/null/v8"
	"go.giteam.ir/giteam/internal/orm"
)

const (
	UserNodeType        NodeType = "User"
	UserTokenNodeType   NodeType = "UserToken"
	UserEmailNodeType   NodeType = "UserEmail"
	UserProfileNodeType NodeType = "UserProfile"
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

// IsNode
func (User) IsNode() {}

// UserFrom Returns an instance of model: `User` from its datasource.
func UserFrom(user *orm.User) *User {
	if user != nil {
		return &User{
			ID:        ToNodeIdentifier(UserNodeType, user.ID),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			RemovedAt: user.RemovedAt,
			IsActive:  user.IsActive,
			IsBanned:  user.IsBanned,
		}
	}

	return nil
}

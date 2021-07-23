package dto

import "go.giteam.ir/giteam/internal/orm"

// UserProfile
type UserProfile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserProfileFrom
func UserProfileFrom(profile *orm.UserProfile) *UserProfile {
	if profile != nil {
		return &UserProfile{
			ID:   ToNodeIdentifier(UserProfileNodeType, profile.ID),
			Name: profile.Name,
		}
	}

	return nil
}

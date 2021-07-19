package dto

// PermissionSubLookup
type PermissionSubLookup string

func (l PermissionSubLookup) String() string {
	switch l {
	case PermissionSubCookieLookup:
		return "COOKIE"
	case PermissionSubHeaderLookup:
		return "HEADER"
	}

	return "INVALID"
}

const (
	PermissionSubHeaderLookup PermissionSubLookup = "HEADER"
	PermissionSubCookieLookup PermissionSubLookup = "COOKIE"
)

type PermissionDefinition struct {
	Obj string `json:"obj"`
	Act string `json:"act"`
}

// PermissionPolicy
type PermissionPolicy struct {
	SubLookup PermissionSubLookup   `json:"subLookup"`
	Def       *PermissionDefinition `json:"def"`
}

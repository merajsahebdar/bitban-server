package conf

// EnvType Represents types for the current running environment.
type EnvType int

const (
	Dev EnvType = iota + 10
	Prod
)

// String Returns the string representation of the current environment.
func (e EnvType) String() string {
	switch e {
	case Prod:
		return "production"
	case Dev:
		return "development"
	}

	return "unknown"
}

// CurrentEnv Keeps the value of current running environment
var CurrentEnv EnvType

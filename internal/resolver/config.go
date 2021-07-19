package resolver

import (
	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/schema"
	"go.uber.org/fx"
)

// ConfigOpt
var ConfigOpt = fx.Provide(newConfig)

// newConfig
func newConfig() schema.Config {
	return schema.Config{
		Directives: schema.DirectiveRoot{
			Guard: (&Guard{
				enforcer: common.GetEnforcerInstance(),
			}).Exec,
		},
	}
}

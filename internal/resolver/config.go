package resolver

import (
	"go.giteam.ir/giteam/internal/component"
	"go.giteam.ir/giteam/internal/controller"
	"go.giteam.ir/giteam/internal/schema"
	"go.uber.org/fx"
)

// ConfigOpt
var ConfigOpt = fx.Provide(newConfig)

// newConfig
func newConfig(accountController *controller.Account) schema.Config {
	return schema.Config{
		Resolvers: &rootResolver{
			validate:          component.GetValidateInstance(),
			accountController: accountController,
		},
	}
}

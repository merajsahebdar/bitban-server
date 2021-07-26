package resolver

import (
	"go.uber.org/fx"
	"regeet.io/api/internal/component"
	"regeet.io/api/internal/controller"
	"regeet.io/api/internal/schema"
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

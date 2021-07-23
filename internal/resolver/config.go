package resolver

import (
	"go.giteam.ir/giteam/internal/common"
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
			validate:          common.GetValidateInstance(),
			accountController: accountController,
		},
	}
}

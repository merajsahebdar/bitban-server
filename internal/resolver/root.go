package resolver

import (
	"github.com/go-playground/validator/v10"
	"regeet.io/api/internal/controller"
	"regeet.io/api/internal/schema"
)

type (
	// rootResolver
	rootResolver struct {
		validate          *validator.Validate
		accountController *controller.Account
	}

	// queryResolver
	queryResolver struct {
		*rootResolver
	}

	// mutationResolver
	mutationResolver struct {
		*rootResolver
	}

	// userResolver
	userResolver struct {
		*rootResolver
	}
)

// Query
func (r *rootResolver) Query() schema.QueryResolver {
	return &queryResolver{
		rootResolver: r,
	}
}

// Mutation
func (r *rootResolver) Mutation() schema.MutationResolver {
	return &mutationResolver{
		rootResolver: r,
	}
}

// User
func (r *rootResolver) User() schema.UserResolver {
	return &userResolver{
		rootResolver: r,
	}
}

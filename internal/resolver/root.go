package resolver

import "go.giteam.ir/giteam/internal/schema"

type (
	// rootResolver
	rootResolver struct{}

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

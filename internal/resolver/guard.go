package resolver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/99designs/gqlgen/graphql"
	"github.com/casbin/casbin/v2"
	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/orm"
)

// Guard
type Guard struct {
	enforcer *casbin.Enforcer
}

// parseObj
func parseObj(obj string, data interface{}) (string, error) {
	t, err := template.New("obj").Parse(obj)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	t.Execute(&res, data)

	return res.String(), nil
}

// Exec
func (d *Guard) Exec(ctx context.Context, parent interface{}, next graphql.Resolver, policy dto.PermissionPolicy) (interface{}, error) {
	var err error

	// Check whether the type is already resolved or not.
	if ret, ok := parent.(*dto.User); ok && ret.ID != "" {
		return next(ctx)
	}

	var user *orm.User
	if user, err = common.AuthorizeUser(ctx, policy.SubLookup); err != nil {
		return nil, AuthenticationErrorFrom(err)
	}

	nextCtx := common.WithAuthorizedUser(ctx, user)

	if policy.Def == nil {
		return next(nextCtx)
	}

	fc := graphql.GetFieldContext(ctx)

	args := fc.Args
	args["parent"] = parent

	var obj string
	if obj, err = parseObj(policy.Def.Obj, args); err != nil {
		panic(err)
	}

	var hasPermission bool
	if hasPermission, err = d.enforcer.Enforce(
		fmt.Sprintf("/users/%d", user.ID),
		common.DefaultUserDomain,
		obj,
		policy.Def.Act,
	); err != nil {
		panic(err)
	}

	if !hasPermission {
		return nil, ForbiddenErrorFrom(
			errors.New("user does not have enough permissions to perform this operation"),
		)
	}

	return next(nextCtx)
}

package common

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"go.uber.org/zap"
)

// DefaultUserDomain
const DefaultUserDomain = "_"

// enforcerLock
var enforcerLock = &sync.Mutex{}

// enforcerInstance
var enforcerInstance *casbin.Enforcer

// GetEnforcerInstance
func GetEnforcerInstance() *casbin.Enforcer {
	if enforcerInstance == nil {
		enforcerLock.Lock()
		defer enforcerLock.Unlock()

		if enforcerInstance == nil {
			enforcerInstance = newEnforcer()
		}
	}

	return enforcerInstance
}

// newEnforcer
func newEnforcer() *casbin.Enforcer {
	var err error

	//
	// Init Model

	var m model.Model
	if m, err = model.NewModelFromString(`
	[request_definition]
	r = sub, dom, obj, act

	[policy_definition]
	p = sub, dom, obj, act

	[role_definition]
	g = _, _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
	`); err != nil {
		Log.Fatal("failed to initialize casbin model", zap.Error(err))
	}

	var a Adapter
	if a, err = newAdapter(); err != nil {
		Log.Fatal("failed to initialize casbin adapter", zap.Error(err))
	}

	//
	// Init Enforcer

	var e *casbin.Enforcer
	if e, err = casbin.NewEnforcer(m, a); err != nil {
		Log.Fatal("failed to initialize casbin enforcer", zap.Error(err))
	}

	e.LoadPolicy()

	return e
}

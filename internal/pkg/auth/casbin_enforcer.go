/*
 * Copyright 2021 Meraj Sahebdar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auth

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"go.uber.org/zap"
	"bitban.io/server/internal/cfg"
)

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
		cfg.Log.Fatal("failed to initialize casbin model", zap.Error(err))
	}

	var a Adapter
	if a, err = newAdapter(); err != nil {
		cfg.Log.Fatal("failed to initialize casbin adapter", zap.Error(err))
	}

	//
	// Init Enforcer

	var e *casbin.Enforcer
	if e, err = casbin.NewEnforcer(m, a); err != nil {
		cfg.Log.Fatal("failed to initialize casbin enforcer", zap.Error(err))
	}

	e.LoadPolicy()

	return e
}

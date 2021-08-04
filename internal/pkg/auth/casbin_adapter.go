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
	"context"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/uptrace/bun"
	"github.com/volatiletech/null/v8"
	"regeet.io/api/internal/pkg/orm"
	"regeet.io/api/internal/pkg/orm/entity"
)

// Adapter
type Adapter interface {
	persist.Adapter
	persist.BatchAdapter
	persist.FilteredAdapter
}

// adapter
type adapter struct {
	db *bun.DB
}

// newAdapter
func newAdapter() (*adapter, error) {
	return &adapter{
		db: orm.GetBunInstance(),
	}, nil
}

// loadPolicyLine
func (a *adapter) loadPolicyLine(m model.Model, line *entity.CasbinRule) {
	arr := []string{line.Ptype, line.V0, line.V1, line.V2, line.V3.String, line.V4.String, line.V5.String}

	var texted string
	if !line.V5.IsZero() {
		texted = strings.Join(arr, ", ")
	} else if !line.V4.IsZero() {
		texted = strings.Join(arr[:6], ", ")
	} else if !line.V3.IsZero() {
		texted = strings.Join(arr[:5], ", ")
	} else if line.V2 != "" {
		texted = strings.Join(arr[:4], ", ")
	} else if line.V1 != "" {
		texted = strings.Join(arr[:3], ", ")
	} else if line.V0 != "" {
		texted = strings.Join(arr[:2], ", ")
	}

	persist.LoadPolicyLine(texted, m)
}

// savePolicyLine
func (*adapter) savePolicyLine(ptype string, rule []string) *entity.CasbinRule {
	line := &entity.CasbinRule{
		Ptype: ptype,
	}

	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = null.StringFrom(rule[3])
	}
	if len(rule) > 4 {
		line.V4 = null.StringFrom(rule[4])
	}
	if len(rule) > 5 {
		line.V5 = null.StringFrom(rule[5])
	}

	return line
}

// IsFiltered
func (*adapter) IsFiltered() bool {
	return true
}

// AddPolicy
func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := a.savePolicyLine(ptype, rule)

	if _, err := a.db.NewInsert().Model(line).Exec(context.Background()); err != nil {
		return err
	}

	return nil
}

// AddPolicies
func (a *adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	var err error

	var tx bun.Tx
	if tx, err = a.db.BeginTx(context.Background(), nil); err != nil {
		return err
	}

	defer func() error {
		return tx.Rollback()
	}()

	for _, rule := range rules {
		line := a.savePolicyLine(ptype, rule)
		if _, err = a.db.NewInsert().Model(line).Exec(context.Background()); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// LoadFilteredPolicy
func (a *adapter) LoadFilteredPolicy(m model.Model, filter interface{}) error {
	return nil
}

// LoadPolicy
func (a *adapter) LoadPolicy(m model.Model) error {
	var lines []*entity.CasbinRule
	if err := a.db.NewSelect().Model(&lines).Scan(context.Background()); err != nil {
		return err
	}

	for _, line := range lines {
		a.loadPolicyLine(m, line)
	}

	return nil
}

// RemoveFilteredPolicy
func (*adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

// RemovePolicy
func (*adapter) RemovePolicy(sec string, ptype string, rules []string) error {
	return nil
}

// RemovePolicies
func (*adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return nil
}

// SavePolicy
func (*adapter) SavePolicy(m model.Model) error {
	return nil
}

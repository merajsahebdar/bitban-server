package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"regeet.io/api/internal/db"
	"regeet.io/api/internal/orm"
)

// Adapter
type Adapter interface {
	persist.Adapter
	persist.BatchAdapter
	persist.FilteredAdapter
}

// adapter
type adapter struct {
	db *sql.DB
}

// newAdapter
func newAdapter() (*adapter, error) {
	return &adapter{
		db: db.GetDbInstance(),
	}, nil
}

// loadPolicyLine
func (a *adapter) loadPolicyLine(m model.Model, p *orm.CasbinRule) {
	line := []string{p.Ptype, p.V0, p.V1, p.V2, p.V3.String, p.V4.String, p.V5.String}

	var texted string
	if !p.V5.IsZero() {
		texted = strings.Join(line, ", ")
	} else if !p.V4.IsZero() {
		texted = strings.Join(line[:6], ", ")
	} else if !p.V3.IsZero() {
		texted = strings.Join(line[:5], ", ")
	} else if p.V2 != "" {
		texted = strings.Join(line[:4], ", ")
	} else if p.V1 != "" {
		texted = strings.Join(line[:3], ", ")
	} else if p.V0 != "" {
		texted = strings.Join(line[:2], ", ")
	}

	persist.LoadPolicyLine(texted, m)
}

// savePolicyLine
func (*adapter) savePolicyLine(ptype string, rule []string) orm.CasbinRule {
	p := orm.CasbinRule{}

	p.Ptype = ptype
	if len(rule) > 0 {
		p.V0 = rule[0]
	}
	if len(rule) > 1 {
		p.V1 = rule[1]
	}
	if len(rule) > 2 {
		p.V2 = rule[2]
	}
	if len(rule) > 3 {
		p.V3 = null.StringFrom(rule[3])
	}
	if len(rule) > 4 {
		p.V4 = null.StringFrom(rule[4])
	}
	if len(rule) > 5 {
		p.V5 = null.StringFrom(rule[5])
	}

	return p
}

// IsFiltered
func (*adapter) IsFiltered() bool {
	return true
}

// AddPolicy
func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	p := a.savePolicyLine(ptype, rule)

	if err := p.Insert(context.Background(), a.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

// AddPolicies
func (a *adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	var err error

	var tx *sql.Tx
	if tx, err = a.db.Begin(); err != nil {
		return err
	}

	defer func() error {
		return tx.Rollback()
	}()

	for _, rule := range rules {
		p := a.savePolicyLine(ptype, rule)
		if err = p.Insert(context.Background(), tx, boil.Infer()); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// LoadFilteredPolicy
func (a *adapter) LoadFilteredPolicy(m model.Model, filter interface{}) error {
	var err error

	var ok bool
	var mods []qm.QueryMod
	if mods, ok = filter.([]qm.QueryMod); !ok {
		return errors.New("invalid adapter filter")
	}

	var policies orm.CasbinRuleSlice
	if policies, err = orm.CasbinRules(mods...).All(context.Background(), a.db); err != nil {
		return err
	}

	for _, p := range policies {
		a.loadPolicyLine(m, p)
	}

	return nil
}

// LoadPolicy
func (a *adapter) LoadPolicy(m model.Model) error {
	var err error

	var policies orm.CasbinRuleSlice
	if policies, err = orm.CasbinRules().All(context.Background(), a.db); err != nil {
		return err
	}

	for _, p := range policies {
		a.loadPolicyLine(m, p)
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

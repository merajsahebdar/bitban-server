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

package validate

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/uptrace/bun"
	"bitban.io/server/internal/cfg"
	"bitban.io/server/internal/pkg/orm"
	"bitban.io/server/internal/pkg/orm/entity"
)

// phoneRegexp
var phoneRegexp *regexp.Regexp

// validateComponentLock
var validateComponentLock = &sync.Mutex{}

// validate
var validateInstance *validator.Validate

// GetValidateInstance
func GetValidateInstance() *validator.Validate {
	if validateInstance == nil {
		validateComponentLock.Lock()
		defer validateComponentLock.Unlock()

		if validateInstance == nil {
			v := validator.New()
			en_translations.RegisterDefaultTranslations(v, cfg.EnTrans)

			//
			// Phone Number Validation

			if compiled, err := regexp.Compile(
				`(0|\+98)?([ ]|-|[()]){0,2}9[1|2|3|4]([ ]|-|[()]){0,2}(?:[0-9]([ ]|-|[()]){0,2}){8}`,
			); err != nil {
				cfg.Log.Fatal(err.Error())
			} else {
				phoneRegexp = compiled
			}

			v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
				number := fl.Field().String()
				if match := phoneRegexp.MatchString(number); match {
					return true
				}

				return false
			})

			v.RegisterTranslation("phone", cfg.EnTrans, func(ut ut.Translator) error {
				return ut.Add("phone", "{0} must be a valid phone number", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				if t, err := ut.T("phone"); err != nil {
					panic(err)
				} else {
					return t
				}
			})

			//
			// Unique Validation

			v.RegisterValidation("notexistsin", func(fl validator.FieldLevel) bool {
				param := strings.Fields(fl.Param())
				datasource := param[0]
				property := param[1]

				switch datasource {
				case "repositories":
					switch property {
					case "address":
						address := fl.Field().String()
						repository := new(entity.Repository)
						if count, err := orm.
							GetBunInstance().
							NewSelect().
							Model(repository).
							Where("? = ?", bun.Ident("repository.address"), address).
							Count(context.Background()); err != nil {
							panic(err)
						} else {
							return count == 0
						}
					}
				case "domains":
					switch property {
					case "address":
						address := fl.Field().String()
						domain := new(entity.Domain)
						if count, err := orm.GetBunInstance().
							NewSelect().
							Model(domain).
							Where("? = ?", bun.Ident("domain.address"), address).
							Count(context.Background()); err != nil {
							panic(err)
						} else {
							return count == 0
						}
					}
				case "emails":
					switch property {
					case "address":
						address := fl.Field().String()
						email := new(entity.Email)
						if count, err := orm.GetBunInstance().
							NewSelect().
							Model(email).
							Where("? = ?", bun.Ident("email.address"), address).
							Count(context.Background()); err != nil {
							panic(err)
						} else {
							return count == 0
						}
					}
				}

				panic(fmt.Errorf("validator is not registered completely"))
			})

			v.RegisterTranslation("notexistsin", cfg.EnTrans, func(ut ut.Translator) error {
				return ut.Add("notexistsin", "{0} must be unique", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				if t, err := ut.T(
					"notexistsin",
					fe.Field(),
					fe.Value().(string),
				); err != nil {
					panic(err)
				} else {
					return t
				}
			})

			validateInstance = v
		}
	}

	return validateInstance
}

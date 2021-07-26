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

package component

import (
	"context"
	"regexp"
	"strings"
	"sync"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/db"
	"regeet.io/api/internal/orm"
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
			en_translations.RegisterDefaultTranslations(v, conf.EnTrans)

			//
			// Phone Number Validation

			if compiled, err := regexp.Compile(
				`(0|\+98)?([ ]|-|[()]){0,2}9[1|2|3|4]([ ]|-|[()]){0,2}(?:[0-9]([ ]|-|[()]){0,2}){8}`,
			); err != nil {
				conf.Log.Fatal(err.Error())
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

			v.RegisterTranslation("phone", conf.EnTrans, func(ut ut.Translator) error {
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

			ctx := context.Background()

			v.RegisterValidation("notexistsin", func(fl validator.FieldLevel) bool {
				param := strings.Fields(fl.Param())
				datasource := param[0]
				property := param[1]

				switch datasource {
				case orm.TableNames.UserEmails:
					switch property {
					case "address":
						address := fl.Field().String()
						var err error
						var exists bool
						if exists, err = orm.UserEmails(
							orm.UserEmailWhere.Address.EQ(address),
						).Exists(ctx, db.GetDbInstance()); err != nil {
							panic(err)
						}

						return !exists
					}
				}

				panic("validator is not registered completely")
			})

			v.RegisterTranslation("notexistsin", conf.EnTrans, func(ut ut.Translator) error {
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

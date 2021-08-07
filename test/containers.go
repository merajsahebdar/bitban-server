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

package test

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// CreatePostgresContainer
func CreatePostgresContainer() {
	if _, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:13",
			ExposedPorts: []string{"5432:5432"},
			Env: map[string]string{
				"POSTGRES_USER":     "regeet",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "regeet",
			},
			WaitingFor: wait.ForSQL("5432", "postgres", func(p nat.Port) string {
				return fmt.Sprintf("postgres://regeet:password@127.0.0.1:%s/regeet?sslmode=disable", p.Port())
			}),
		},
		Started: true,
	}); err != nil {
		log.Fatal(err.Error())
	}
}

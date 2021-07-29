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

package cfg

// GitBackend
type GitBackend string

const (
	GitBackendBin GitBackend = "bin"
	GitBackendGo  GitBackend = "go"
)

// String
func (b GitBackend) String() string {
	switch b {
	case GitBackendBin:
		return "bin"
	case GitBackendGo:
		return "go"
	}

	panic("git backend not implemented")
}

// Cog
var Cog struct {
	App struct {
		Host string `yaml:"host" default:"0.0.0.0"`
		Port int    `yaml:"port" default:"8080"`
	} `yaml:"app"`
	Git struct {
		Backend GitBackend `yaml:"backend"`
	} `yaml:"git"`
	Security struct {
		AccessTokenExpiresAt  int `yaml:"accessTokenExpiresAt" default:"60"`
		RefreshTokenExpiresAt int `yaml:"refreshTokenExpiresAt" default:"259200"`
	} `yaml:"security"`
	Database struct {
		Host   string `yaml:"host" default:"127.0.0.1"`
		Port   int    `yaml:"port" default:"5432"`
		Dbname string `yaml:"dbname"`
		User   string `yaml:"user"`
		Pass   string `yaml:"pass"`
	} `yaml:"database"`
	Redis struct {
		Url string `yaml:"url"`
	}
	Ssh struct {
		Key struct {
			PublicKey  string `yaml:"publicKey"`
			PrivateKey string `yaml:"privateKey"`
			Passphrase string `yaml:"passphrase"`
		} `yaml:"key"`
	} `yaml:"ssh"`
	Jwt struct {
		Key struct {
			PublicKey  string `yaml:"publicKey"`
			PrivateKey string `yaml:"privateKey"`
			Passphrase string `yaml:"passphrase"`
		} `yaml:"key"`
	} `yaml:"jwt"`
}

// IsGoBacked
func IsGoBackend() bool {
	return Cog.Git.Backend == GitBackendGo
}

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

// GitStorage
type GitStorage string

const (
	GitStorageMem GitStorage = "mem"
	GitStorageFs  GitStorage = "fs"
)

// Cog
var Cog struct {
	App struct {
		Host string `yaml:"host" default:"0.0.0.0"`
		Port int    `yaml:"port" default:"8080"`
	} `yaml:"app"`
	Git struct {
		Backend GitBackend `yaml:"backend"`
		Storage GitStorage `yaml:"storage"`
		Configs struct {
			Init struct {
				DefaultBranch string `yaml:"defaultBranch" default:"main"`
			} `yaml:"init"`
		} `yaml:"configs"`
	} `yaml:"git"`
	Security struct {
		AccessTokenExpiresAt  int `yaml:"accessTokenExpiresAt" default:"60"`
		RefreshTokenExpiresAt int `yaml:"refreshTokenExpiresAt" default:"259200"`
	} `yaml:"security"`
	Database struct {
		Host   string `yaml:"host" default:"127.0.0.1"`
		Port   int    `yaml:"port" default:"5432"`
		Dbname string `yaml:"dbname" default:"bitban"`
		User   string `yaml:"user" default:"bitban"`
		Pass   string `yaml:"pass" default:"password"`
	} `yaml:"database"`
	Redis struct {
		Url string `yaml:"url" default:"redis://127.0.0.1:6379/0"`
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

package conf

// Cog
var Cog struct {
	App struct {
		Host string `yaml:"host" default:"0.0.0.0"`
		Port int    `yaml:"port" default:"8080"`
	} `yaml:"app"`
	Storage struct {
		Dir string `yaml:"dir"`
	} `yaml:"storage"`
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
	Amqp struct {
		Uri string `yaml:"uri"`
	} `yaml:"amqp"`
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
		PublicKey  string `yaml:"publicKey"`
		PrivateKey string `yaml:"privateKey"`
		Passphrase string `yaml:"passphrase"`
	} `yaml:"jwt"`
}

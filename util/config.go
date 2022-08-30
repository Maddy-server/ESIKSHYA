package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		User     string `toml:"user"`
		Password string `toml:"password"`
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Name     string `toml:"name"`
	} `toml:"database"`
	SMS struct {
		Token string `toml:"token"`
		From  string `toml:"from"`
	} `toml:"sms"`
	Logging struct {
		Level string `toml:"logging"`
	} `toml:"logging"`

	Server struct {
		Listen string `toml:"listen"`
	} `toml:"server"`

	Firebase struct {
		CredentialPath string `toml:"credentialPath"`
	} `toml:"firebase"`

	AWS struct {
		BucketName string `toml:"bucketName"`
		Region     string `toml:"region"`
		Secret     string `toml:"secret"`
		ID         string `toml:"id"`
	} `toml:"aws"`

	Mail struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"mail"`

	Keys struct {
		PrivateKey string `toml:"privateKey"`
		PublicKey  string `toml:"publicKey"`
	} `toml:"keys"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

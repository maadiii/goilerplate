package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type RefreshToken struct {
	Secret    []byte
	Algorithm string
	MaxAge    uint
	Secure    bool
	HTTPOnly  bool
	Path      string
}

type JsonWebToken struct {
	Secret       []byte
	Algorithm    string
	MaxAge       uint
	HTTPOnly     bool
	RefreshToken RefreshToken
}

type TLSConf struct {
	CRT string
	Key string
}

type Config struct {
	SecretKey      []byte
	BlockSecretKey []byte
	JWT            JsonWebToken
	Static         string
	TLS            TLSConf
}

var debugConfig = &Config{
	// It is recommended to use a key with 32 or 64 bytes.
	SecretKey: []byte(strings.Repeat("#", 32)),
	// valid length are 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	BlockSecretKey: []byte(strings.Repeat("#", 32)),
	JWT: JsonWebToken{
		Secret:    []byte("<JWT-SECRET>"),
		Algorithm: "HS256",
		MaxAge:    uint(900),
		HTTPOnly:  true,
		RefreshToken: RefreshToken{
			Secret:    []byte("<JWT-REFRESH-SECRET>"),
			Algorithm: "HS256",
			MaxAge:    uint(86400),
			Secure:    true,
			HTTPOnly:  true,
			Path:      "/",
		},
	},
	Static: "./static",
	TLS: TLSConf{
		CRT: "./server.crt",
		Key: "./server.key",
	},
}

func InitConfig() (*Config, error) {
	var config *Config

	if viper.ConfigFileUsed() == EMPTY {
		config = debugConfig
	} else {
		config = &Config{
			SecretKey:      []byte(viper.GetString(SECRET_KEY)),
			BlockSecretKey: []byte(viper.GetString(BLOCK_SECRET_KEY)),
			JWT: JsonWebToken{
				Secret:    []byte(viper.GetString(SECRET_KEY)),
				Algorithm: viper.GetString(JWT_ALGORITHM),
				MaxAge:    viper.GetUint(JWT_MAXAGE),
				HTTPOnly:  viper.GetBool(JWT_HTTPONLY),
				RefreshToken: RefreshToken{
					Secret:    []byte(viper.GetString(REFRESH_TOKEN_SECRET)),
					Algorithm: viper.GetString(REFRESH_TOKEN_ALGORITHM),
					MaxAge:    viper.GetUint(REFRESH_TOKEN_MAXAGE),
					Secure:    viper.GetBool(REFRESH_TOKEN_SECURE),
					HTTPOnly:  viper.GetBool(REFRESH_TOKEN_HTTPONLY),
					Path:      viper.GetString(REFRESH_TOKEN_PATH),
				},
			},
			Static: viper.GetString(STATIC),
			TLS: TLSConf{
				CRT: viper.GetString(CRT),
				Key: viper.GetString(KEY),
			},
		}
	}

	preError := "app/config.go InitConfig(),"
	postError := "is not set in config file"

	if len(config.SecretKey) < 16 {
		return nil, errors.New(
			fmt.Sprintf("%s %s %s", preError, SECRET_KEY, postError),
		)
	}

	if len(config.JWT.Secret) < 8 {
		return nil, errors.New(
			fmt.Sprintf(
				"%s %s %s %s",
				preError,
				JWT_SECRET,
				postError,
				"or lesser than 16 byte",
			),
		)
	}

	if config.JWT.Algorithm != HS256 &&
		config.JWT.Algorithm != HS384 &&
		config.JWT.Algorithm != HS512 {
		return nil, errors.New(
			fmt.Sprintf(
				"%s %s %s %s",
				preError,
				JWT_ALGORITHM,
				postError,
				"or not in (HS256, HS384, HS512)",
			),
		)
	}

	if config.JWT.MaxAge == 0 {
		return nil, errors.New(
			fmt.Sprintf("%s %s %s", preError, JWT_MAXAGE, postError),
		)
	}

	if len(config.JWT.RefreshToken.Secret) < 8 {
		return nil, errors.New(
			fmt.Sprintf(
				"%s %s %s %s",
				preError,
				REFRESH_TOKEN_SECRET,
				postError,
				"or lesser than 16 byte",
			),
		)
	}

	if config.JWT.RefreshToken.Algorithm != HS256 &&
		config.JWT.Algorithm != HS384 &&
		config.JWT.Algorithm != HS512 {
		return nil, errors.New(
			fmt.Sprintf(
				"%s %s %s %s",
				preError,
				REFRESH_TOKEN_ALGORITHM,
				postError,
				"or not in (HS256, HS384, HS512",
			),
		)
	}

	if config.JWT.RefreshToken.MaxAge == 0 {
		return nil, errors.New(
			fmt.Sprintf(
				"%s %s %s", preError, REFRESH_TOKEN_MAXAGE, postError,
			),
		)
	}

	if len(config.JWT.RefreshToken.Path) == 0 {
		return nil, errors.New(
			fmt.Sprintf("%s %s %s", preError, REFRESH_TOKEN_PATH, postError),
		)
	}

	if len(config.TLS.Key) == 0 {
		return nil, errors.New(
			fmt.Sprintf("%s %s %s", preError, KEY, postError),
		)
	}

	if len(config.TLS.CRT) == 0 {
		return nil, errors.New(
			fmt.Sprintf("%s %s %s", preError, CRT, postError),
		)
	}

	return config, nil
}

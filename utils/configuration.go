package utils

import (
	"encoding/json"
	"os"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Configuration struct {
	HttpAddress             string                 `json:"HttpAddress"`
	InvitationCode          string                 `json:"InvitationCode"`
	AdminPassword           string                 `json:"AdminPassword"`
	AuthorizationCodeLength int                    `json:"AuthorizationCodeLength"`
	RefreshTokenLength      int                    `json:"RefreshTokenLength"`
	PasswordMinLength       int                    `json:"PasswordMinLength"`
	LogFilePath             string                 `json:"LogFilePath"`
	JwtIssuer               string                 `json:"JwtIssuer"`
	JwtPrivateKey           map[string]interface{} `json:"JwtPrivateKey"`
	JwtPublicKey            map[string]interface{} `json:"JwtPublicKey"`

	ParsedJwtPrivateKey jwk.Key `json:"-"`
	ParsedJwtPublicKey  jwk.Key `json:"-"`
}

func (c *Configuration) Init(mode string) error {
	filename := "conf-" + mode + ".json"
	confJson, err := os.ReadFile("./" + filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(confJson, &c); err != nil {
		return err
	}

	privateKeyByte, err := json.Marshal(c.JwtPrivateKey)
	if err != nil {
		return err
	}
	c.ParsedJwtPrivateKey, err = jwk.ParseKey(privateKeyByte)
	if err != nil {
		return err
	}
	c.ParsedJwtPublicKey, err = jwk.PublicKeyOf(c.ParsedJwtPrivateKey)
	if err != nil {
		return err
	}

	return nil
}

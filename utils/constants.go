package utils

import (
	"encoding/json"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"os"
)

type Configuration struct {
	HttpAddress             string                 `json:"HttpAddress"`
	InvitationCode          string                 `json:"InvitationCode"`
	AdminPassword           string                 `json:"AdminPassword"`
	AuthorizationCodeLength int                    `json:"AuthorizationCodeLength"`
	RefreshTokenLength      int                    `json:"RefreshTokenLength"`
	JwtIssuer               string                 `json:"JwtIssuer"`
	JwtPrivateKey           map[string]interface{} `json:"JwtPrivateKey"`
	JwtPublicKey            map[string]interface{} `json:"JwtPublicKey"`

	ParsedJwtPrivateKey jwk.Key `json:"-"`
	ParsedJwtPublicKey  jwk.Key `json:"-"`
}

var Conf Configuration

func Init(mode string) error {
	filename := "conf-" + mode + ".json"
	confJson, err := os.ReadFile("./" + filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(confJson, &Conf); err != nil {
		return err
	}

	AuthorizationCodeMap = make(map[string]AuthorizationInfo, 0)

	privateKeyByte, err := json.Marshal(Conf.JwtPrivateKey)
	if err != nil {
		return err
	}
	Conf.ParsedJwtPrivateKey, err = jwk.ParseKey(privateKeyByte)
	if err != nil {
		return err
	}
	Conf.ParsedJwtPublicKey, err = jwk.PublicKeyOf(Conf.ParsedJwtPrivateKey)
	if err != nil {
		return err
	}

	return nil
}

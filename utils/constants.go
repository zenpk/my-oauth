package utils

import (
	"encoding/json"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"os"
)

type Configuration struct {
	HttpAddress             string
	InvitationCode          string
	AdminPassword           string
	AuthorizationCodeLength int
	RefreshTokenLength      int
	JwtIssuer               string
	JwtPrivateKey           string
	JwtPublicKey            string

	parsedJwtPrivateKey jwk.Key
	parsedJwtPublicKey  jwk.Key
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

	Conf.parsedJwtPrivateKey, err = jwk.ParseKey([]byte(Conf.JwtPrivateKey))
	if err != nil {
		return err
	}
	Conf.parsedJwtPublicKey, err = jwk.PublicKeyOf(Conf.parsedJwtPrivateKey)
	if err != nil {
		return err
	}

	return nil
}

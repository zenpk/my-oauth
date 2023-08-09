package utils

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	HttpAddress             string
	InvitationCode          string
	AdminPassword           string
	AuthorizationCodeLength int
	JwtIssuer               string
	JwtPrivateKey           string
	JwtPublicKey            string
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

	return nil
}

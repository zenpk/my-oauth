package util

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	HttpAddress             string `json:"httpAddress"`
	InvitationCode          string `json:"invitationCode"`
	AdminPassword           string `json:"adminPassword"`
	AuthorizationCodeLength int    `json:"authorizationCodeLength"`
	RefreshTokenLength      int    `json:"refreshTokenLength"`
	PasswordMinLength       int    `json:"passwordMinLength"`
	JwtIssuer               string `json:"jwtIssuer"`
	DbPath                  string `json:"dbPath"`
	LogFilePath             string `json:"logFilePath"`
	RsaPrivateKeyPath       string `json:"rsaPrivateKeyPath"`
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
	return nil
}

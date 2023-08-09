package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

type AuthorizationInfo struct {
	ClientId      string
	Uuid          string
	CodeChallenge string
}

var AuthorizationCodeMap map[string]AuthorizationInfo

func GenAuthorizationCode(info AuthorizationInfo) (string, error) {
	code, err := RandString(Conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}
	AuthorizationCodeMap[code] = info
	return code, nil
}

func VerifyAuthorizationCode(code string, codeVerifier string) bool {
	info, ok := AuthorizationCodeMap[code]
	if !ok {
		return false
	}
	checksum := sha256.Sum256([]byte(codeVerifier))
	// use base64.RawURLEncoding to omit padding
	return base64.RawURLEncoding.EncodeToString(checksum[:]) == info.CodeChallenge
}

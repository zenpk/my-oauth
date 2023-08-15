package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type AuthorizationInfo struct {
	ClientId      string
	Uuid          string
	Username      string
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

func VerifyAuthorizationCode(code string, codeVerifier string) (AuthorizationInfo, error) {
	info, ok := AuthorizationCodeMap[code]
	if !ok {
		return AuthorizationInfo{}, errors.New("invalid authorization code")
	}
	checksum := sha256.Sum256([]byte(codeVerifier))
	// use base64.RawURLEncoding to omit padding
	match := base64.RawURLEncoding.EncodeToString(checksum[:]) == info.CodeChallenge
	if !match {
		return AuthorizationInfo{}, errors.New("code challenge failed")
	}
	delete(AuthorizationCodeMap, code) // one-time usage
	return info, nil
}

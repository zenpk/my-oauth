package util

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type AuthorizationInfo struct {
	ClientId             int64
	UserId               int64
	CodeChallenge        string
	Context              string
	conf                 *Configuration
	authorizationCodeMap map[string]*AuthorizationInfo
}

func (a *AuthorizationInfo) Init(conf *Configuration) {
	a.conf = conf
	a.authorizationCodeMap = make(map[string]*AuthorizationInfo)
}

func (a *AuthorizationInfo) GenAuthorizationCode(info *AuthorizationInfo) (string, error) {
	code, err := RandString(a.conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}
	a.authorizationCodeMap[code] = info
	return code, nil
}

// VerifyAuthorizationCode verifier(base64url) -> []byte -> sha256([]byte) -> base64url should == challenge(base64url)
func (a *AuthorizationInfo) VerifyAuthorizationCode(code string, codeVerifier string) (*AuthorizationInfo, error) {
	info, ok := a.authorizationCodeMap[code]
	if !ok {
		return nil, errors.New("invalid authorization code")
	}
	checksum := sha256.Sum256([]byte(codeVerifier))
	// use base64.RawURLEncoding to omit padding
	match := base64.RawURLEncoding.EncodeToString(checksum[:]) == info.CodeChallenge
	if !match {
		return nil, errors.New("code challenge failed")
	}
	delete(a.authorizationCodeMap, code) // one-time usage
	return info, nil
}

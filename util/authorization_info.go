package util

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

const authorizationCodeTTL = 5 * time.Minute

type authCodeEntry struct {
	info      *AuthorizationInfo
	expiresAt time.Time
}

type AuthorizationInfo struct {
	ClientId             int64
	UserId               int64
	CodeChallenge        string
	Context              string
	conf                 *Configuration
	mu                   sync.RWMutex
	authorizationCodeMap map[string]*authCodeEntry
}

func (a *AuthorizationInfo) Init(conf *Configuration) {
	a.conf = conf
	a.authorizationCodeMap = make(map[string]*authCodeEntry)
	go a.cleanupLoop()
}

func (a *AuthorizationInfo) GenAuthorizationCode(info *AuthorizationInfo) (string, error) {
	code, err := RandString(a.conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}
	a.mu.Lock()
	a.authorizationCodeMap[code] = &authCodeEntry{
		info:      info,
		expiresAt: time.Now().Add(authorizationCodeTTL),
	}
	a.mu.Unlock()
	return code, nil
}

func (a *AuthorizationInfo) VerifyAuthorizationCode(code string, codeVerifier string) (*AuthorizationInfo, error) {
	a.mu.Lock()
	entry, ok := a.authorizationCodeMap[code]
	if !ok {
		a.mu.Unlock()
		return nil, errors.New("invalid authorization code")
	}
	delete(a.authorizationCodeMap, code)
	a.mu.Unlock()

	if time.Now().After(entry.expiresAt) {
		return nil, errors.New("authorization code expired")
	}
	checksum := sha256.Sum256([]byte(codeVerifier))
	match := base64.RawURLEncoding.EncodeToString(checksum[:]) == entry.info.CodeChallenge
	if !match {
		return nil, errors.New("code challenge failed")
	}
	return entry.info, nil
}

func (a *AuthorizationInfo) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		a.mu.Lock()
		for code, entry := range a.authorizationCodeMap {
			if now.After(entry.expiresAt) {
				delete(a.authorizationCodeMap, code)
			}
		}
		a.mu.Unlock()
	}
}

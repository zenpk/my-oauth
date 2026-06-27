package util

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

const authorizationCodeTTL = 5 * time.Minute
const authorizationCodeCleanupInterval = 1 * time.Minute

type authCodeEntry struct {
	info      *AuthorizationInfo
	expiresAt time.Time
}

type AuthorizationInfo struct {
	ClientId             int64
	UserId               int64
	RedirectUri          string
	Scope                string
	State                string
	Nonce                string
	CodeChallenge        string
	conf                 *Configuration
	mu                   sync.RWMutex
	authorizationCodeMap map[string]*authCodeEntry
	lastCleanup          time.Time
}

func (a *AuthorizationInfo) Init(conf *Configuration) {
	a.conf = conf
	a.authorizationCodeMap = make(map[string]*authCodeEntry)
	a.lastCleanup = time.Now()
}

func (a *AuthorizationInfo) GenAuthorizationCode(info *AuthorizationInfo) (string, error) {
	code, err := RandString(a.conf.AuthorizationCodeLength)
	if err != nil {
		return "", err
	}
	now := time.Now()
	a.mu.Lock()
	a.cleanupExpiredLocked(now)
	a.authorizationCodeMap[code] = &authCodeEntry{
		info:      info,
		expiresAt: now.Add(authorizationCodeTTL),
	}
	a.mu.Unlock()
	return code, nil
}

func (a *AuthorizationInfo) VerifyAuthorizationCode(code string, codeVerifier string) (*AuthorizationInfo, error) {
	if len(codeVerifier) < 43 || len(codeVerifier) > 128 {
		return nil, errors.New("invalid code verifier")
	}
	now := time.Now()
	a.mu.Lock()
	a.cleanupExpiredLocked(now)
	entry, ok := a.authorizationCodeMap[code]
	if !ok {
		a.mu.Unlock()
		return nil, errors.New("invalid authorization code")
	}
	delete(a.authorizationCodeMap, code)
	a.mu.Unlock()

	if now.After(entry.expiresAt) {
		return nil, errors.New("authorization code expired")
	}
	checksum := sha256.Sum256([]byte(codeVerifier))
	match := base64.RawURLEncoding.EncodeToString(checksum[:]) == entry.info.CodeChallenge
	if !match {
		return nil, errors.New("code challenge failed")
	}
	return entry.info, nil
}

func (a *AuthorizationInfo) cleanupExpiredLocked(now time.Time) {
	if now.Sub(a.lastCleanup) < authorizationCodeCleanupInterval {
		return
	}
	for code, entry := range a.authorizationCodeMap {
		if now.After(entry.expiresAt) {
			delete(a.authorizationCodeMap, code)
		}
	}
	a.lastCleanup = now
}

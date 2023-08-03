package utils

import (
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

type Payload struct {
	Uuid     string
	Username string
	ClientId string
}

func GenerateJwt(payload Payload) (string, error) {
	token, err := jwt.NewBuilder().
		Audience([]string{payload.ClientId}).
		IssuedAt(time.Now()).
		Issuer(Issuer).
		Expiration(time.Now().Add(AccessTokenAge*time.Hour)).
		NotBefore(time.Now()).
		Claim("uuid", payload.Uuid).
		Claim("username", payload.Username).
		Build()
	if err != nil {
		return "", err
	}
	rawKey := []byte(JwtPrivateKey)
	jwkKey, err := jwk.FromRaw(rawKey)
	if err != nil {
		return "", err
	}
	serialized, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, jwkKey))
	if err != nil {
		return "", err
	}
	return string(serialized), nil
}

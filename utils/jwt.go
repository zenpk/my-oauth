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
}

func GenerateJwt(payload Payload, clientId string, tokenAge time.Duration) (string, error) {
	token, err := jwt.NewBuilder().
		Audience([]string{clientId}).
		IssuedAt(time.Now()).
		Issuer(Conf.JwtIssuer).
		Expiration(time.Now().Add(tokenAge*time.Hour)).
		NotBefore(time.Now()).
		Claim("uuid", payload.Uuid).
		Claim("username", payload.Username).
		Build()
	if err != nil {
		return "", err
	}
	priKey, err := jwk.ParseKey([]byte(Conf.JwtPrivateKey))
	if err != nil {
		return "", err
	}
	serialized, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, priKey))
	if err != nil {
		return "", err
	}
	return string(serialized), nil
}

package utils

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/zenpk/my-oauth/db"
	"time"
)

type Payload struct {
	Uuid     string
	Username string
	ClientId string
}

func GenerateJwt(payload Payload, tokenAge time.Duration) (string, error) {
	token, err := jwt.NewBuilder().
		Audience([]string{payload.ClientId}).
		IssuedAt(time.Now()).
		Issuer(Conf.JwtIssuer).
		Expiration(time.Now().Add(tokenAge)).
		NotBefore(time.Now()).
		Claim("uuid", payload.Uuid).
		Claim("username", payload.Username).
		Claim("clientId", payload.ClientId).
		Build()
	if err != nil {
		return "", err
	}
	serialized, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, Conf.ParsedJwtPrivateKey))
	if err != nil {
		return "", err
	}
	return string(serialized), nil
}

func VerifyJwt(token string) error {
	_, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256, Conf.ParsedJwtPublicKey))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to verify JWS: %s\n", err))
	}
	return nil
}

func GenAndInsertRefreshToken(payload Payload, tokenAge time.Duration) (string, error) {
	refreshToken, err := RandString(Conf.RefreshTokenLength)
	if err != nil {
		return "", err
	}
	if err := db.TableRefreshToken.Insert(db.RefreshToken{
		Token:      refreshToken,
		ClientId:   payload.ClientId,
		Uuid:       payload.Uuid,
		Username:   payload.Username,
		ExpireTime: time.Now().Add(tokenAge),
	}); err != nil {
		return "", err
	}
	return refreshToken, nil
}

func GetAndCleanRefreshToken(refreshToken string) (db.RefreshToken, error) {
	tokens, err := db.TableRefreshToken.All()
	if err != nil {
		return db.RefreshToken{}, err
	}
	for _, token := range tokens {
		// delete expired
		if token.(db.RefreshToken).ExpireTime.After(time.Now()) {
			if err := db.TableRefreshToken.Delete(db.RefreshTokenToken, token.(db.RefreshToken).Token); err != nil {
				return db.RefreshToken{}, err
			}
			continue
		}
		if token.(db.RefreshToken).Token == refreshToken {
			// also delete this, because a new one will be generated
			if err := db.TableRefreshToken.Delete(db.RefreshTokenToken, token.(db.RefreshToken).Token); err != nil {
				return db.RefreshToken{}, err
			}
			return token.(db.RefreshToken), nil
		}
	}
	return db.RefreshToken{}, errors.New("no valid refresh token found")
}

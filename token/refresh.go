package token

import (
	"errors"
	"time"

	"github.com/zenpk/my-oauth/db"
)

type RefreshToken struct{}

func (r *RefreshToken) GenAndInsertRefreshToken(dbInstance *db.Db, payload Payload, tokenAge time.Duration) (string, error) {
	refreshToken, err := RandString(Conf.RefreshTokenLength)
	if err != nil {
		return "", err
	}
	if err := dbInstance.TableRefreshToken.Insert(db.RefreshToken{
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

func (r *RefreshToken) GetAndCleanRefreshToken(dbInstance *db.Db, refreshToken string) (db.RefreshToken, error) {
	tokens, err := dbInstance.TableRefreshToken.All()
	if err != nil {
		return db.RefreshToken{}, err
	}
	for _, token := range tokens {
		// delete expired
		if token.(db.RefreshToken).ExpireTime.Before(time.Now()) {
			if err := dbInstance.TableRefreshToken.Delete(db.RefreshTokenToken, token.(db.RefreshToken).Token); err != nil {
				return db.RefreshToken{}, err
			}
			continue
		}
		if token.(db.RefreshToken).Token == refreshToken {
			return token.(db.RefreshToken), nil
		}
	}
	return db.RefreshToken{}, errors.New("no valid refresh token found")
}

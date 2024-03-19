package service

import (
	"time"

	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/token"
	"github.com/zenpk/my-oauth/util"
)

func (s *Service) GenAndInsertRefreshToken(claims *token.Claims, client *dal.Client, user *dal.User) (string, error) {
	refreshToken, err := util.RandString(s.conf.RefreshTokenLength)
	if err != nil {
		return "", err
	}
	expireTime := time.Now().Add(time.Duration(client.RefreshTokenAge) * time.Hour)
	if err := s.db.RefreshTokens.Insert(&dal.RefreshToken{
		Token:      refreshToken,
		ClientId:   client.Id,
		UserId:     user.Id,
		ExpireTime: &expireTime,
	}); err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (s *Service) CleanAndGetRefreshToken(token string) (*dal.RefreshToken, error) {
	// clean expired first so that the selected result is valid
	if err := s.db.RefreshTokens.CleanExpired(); err != nil {
		return nil, err
	}
	refreshToken, err := s.db.RefreshTokens.SelectByToken(token)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

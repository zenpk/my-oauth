package service

import (
	"log"
	"time"

	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/token"
	"github.com/zenpk/my-oauth/util"
)

func (s *Service) GenAndInsertRefreshToken(claims *token.Claims, client *dal.Client, user *dal.User) (string, error) {
	// invalidate old refresh tokens for this user+client pair
	if err := s.db.RefreshTokens.DeleteByUserAndClient(user.Id, client.Id); err != nil {
		return "", err
	}
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

func (s *Service) StartCleanupJob(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.db.RefreshTokens.CleanExpired(); err != nil {
				log.Printf("cleanup job error: %v\n", err)
			}
		}
	}()
}

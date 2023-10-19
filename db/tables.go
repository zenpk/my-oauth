package db

import (
	scd "github.com/zenpk/safe-csv-db"
	"strconv"
	"time"
)

const (
	UserUuid     = 0
	UserUsername = 1
	UserPassword = 2

	ClientId              = 0
	ClientSecret          = 1
	ClientRedirects       = 2
	ClientAccessTokenAge  = 3
	ClientRefreshTokenAge = 4

	RefreshTokenToken      = 0
	RefreshTokenClientId   = 1
	RefreshTokenUuid       = 2
	RefreshTokenUsername   = 3
	RefreshTokenExpireTime = 4
)

type User struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) ToRow() ([]string, error) {
	row := make([]string, 3)
	row[UserUuid] = u.Uuid
	row[UserUsername] = u.Username
	row[UserPassword] = u.Password
	return row, nil
}

func (u User) FromRow(row []string) (scd.RecordType, error) {
	return User{
		Uuid:     row[UserUuid],
		Username: row[UserUsername],
		Password: row[UserPassword],
	}, nil
}

type Client struct {
	Id              string `json:"id"`
	Secret          string `json:"secret"`
	Redirects       string `json:"redirects"`
	AccessTokenAge  int    `json:"accessTokenAge"`  // hour
	RefreshTokenAge int    `json:"refreshTokenAge"` // hour
}

func (c Client) ToRow() ([]string, error) {
	row := make([]string, 5)
	row[ClientId] = c.Id
	row[ClientSecret] = c.Secret
	row[ClientRedirects] = c.Redirects
	row[ClientAccessTokenAge] = strconv.Itoa(c.AccessTokenAge)
	row[ClientRefreshTokenAge] = strconv.Itoa(c.RefreshTokenAge)
	return row, nil
}

func (c Client) FromRow(row []string) (scd.RecordType, error) {
	var err error
	var newClient Client
	newClient.Id = row[ClientId]
	newClient.Secret = row[ClientSecret]
	newClient.Redirects = row[ClientRedirects]
	newClient.AccessTokenAge, err = strconv.Atoi(row[ClientAccessTokenAge])
	if err != nil {
		return nil, err
	}
	newClient.RefreshTokenAge, err = strconv.Atoi(row[ClientRefreshTokenAge])
	if err != nil {
		return nil, err
	}
	return newClient, nil
}

type RefreshToken struct {
	Token      string    `json:"string"`
	ClientId   string    `json:"clientId"`
	Uuid       string    `json:"uuid"`
	Username   string    `json:"username"`
	ExpireTime time.Time `json:"expireTime"` // UNIX second
}

func (r RefreshToken) ToRow() ([]string, error) {
	row := make([]string, 5)
	row[RefreshTokenToken] = r.Token
	row[RefreshTokenClientId] = r.ClientId
	row[RefreshTokenUuid] = r.Uuid
	row[RefreshTokenUsername] = r.Username
	row[RefreshTokenExpireTime] = strconv.FormatInt(r.ExpireTime.Unix(), 10)
	return row, nil
}

func (r RefreshToken) FromRow(row []string) (scd.RecordType, error) {
	var newRefreshToken RefreshToken
	newRefreshToken.Token = row[RefreshTokenToken]
	newRefreshToken.ClientId = row[RefreshTokenClientId]
	newRefreshToken.Uuid = row[RefreshTokenUuid]
	newRefreshToken.Username = row[RefreshTokenUsername]
	secondTime, err := strconv.ParseInt(row[RefreshTokenExpireTime], 10, 64)
	if err != nil {
		return nil, err
	}
	newRefreshToken.ExpireTime = time.UnixMilli(secondTime * 1000)
	return newRefreshToken, nil
}

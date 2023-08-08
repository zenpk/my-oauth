package db

import (
	scd "github.com/zenpk/safe-csv-db"
	"strconv"
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

	RefreshTokenClientId   = 0
	RefreshTokenToken      = 1
	RefreshTokenExpireTime = 2
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

func (u User) FromRow(row []string) (scd.Table, error) {
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
	RefreshTokenAge int    `json:"RefreshTokenAge"` // hour
}

func (c Client) ToRow() ([]string, error) {
	row := make([]string, 4)
	row[ClientId] = c.Id
	row[ClientSecret] = c.Secret
	row[ClientRedirects] = c.Redirects
	row[ClientAccessTokenAge] = strconv.Itoa(c.AccessTokenAge)
	row[ClientRefreshTokenAge] = strconv.Itoa(c.RefreshTokenAge)
	return row, nil
}

func (c Client) FromRow(row []string) (scd.Table, error) {
	var err error
	var newClient Client
	newClient.Id = row[ClientId]
	newClient.Secret = row[ClientSecret]
	newClient.Redirects = row[ClientRedirects]
	newClient.AccessTokenAge, err = strconv.Atoi(row[ClientAccessTokenAge])
	if err != nil {
		return nil, err
	}
	newClient.RefreshTokenAge, err = strconv.Atoi(row[ClientAccessTokenAge])
	if err != nil {
		return nil, err
	}
	return newClient, nil
}

type RefreshToken struct {
	ClientId   string `json:"clientId"`
	Token      string `json:"string"`
	ExpireTime int64  `json:"expireTime"` // UNIX ms
}

func (r RefreshToken) ToRow() ([]string, error) {
	row := make([]string, 3)
	row[RefreshTokenClientId] = r.ClientId
	row[RefreshTokenToken] = r.Token
	row[RefreshTokenExpireTime] = strconv.FormatInt(r.ExpireTime, 10)
	return row, nil
}

func (r RefreshToken) FromRow(row []string) (scd.Table, error) {
	var err error
	var newRefreshToken RefreshToken
	newRefreshToken.ClientId = row[RefreshTokenClientId]
	newRefreshToken.Token = row[RefreshTokenToken]
	newRefreshToken.ExpireTime, err = strconv.ParseInt(row[RefreshTokenExpireTime], 10, 64)
	if err != nil {
		return nil, err
	}
	return newRefreshToken, nil
}

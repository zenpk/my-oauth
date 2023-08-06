package db

import (
	"strconv"
)

const (
	UserUuid     = 0
	UserUsername = 1
	UserPassword = 2

	ClientId              = 0
	ClientSecret          = 1
	ClientRedirect        = 2
	ClientOwner           = 3
	ClientAccessTokenAge  = 4
	ClientRefreshTokenAge = 5

	RefreshTokenClientId   = 0
	RefreshTokenToken      = 1
	RefreshTokenExpireTime = 2
)

type User struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) StructToRow(user User) []string {
	row := make([]string, 3)
	row[UserUuid] = user.Uuid
	row[UserUsername] = u.Username
	row[UserPassword] = u.Password
	return row
}

type Client struct {
	Id              string `json:"id"`
	Secret          string `json:"secret"`
	Redirects       string `json:"redirects"`
	Owner           string `json:"owner"`
	AccessTokenAge  int    `json:"accessTokenAge"`  // hour
	RefreshTokenAge int    `json:"RefreshTokenAge"` // hour
}

func (c Client) StructToRow(client Client) []string {
	row := make([]string, 4)
	row[ClientId] = client.Id
	row[ClientSecret] = client.Secret
	row[ClientRedirect] = client.Redirects
	row[ClientOwner] = client.Owner
	row[ClientAccessTokenAge] = strconv.Itoa(client.AccessTokenAge)
	row[ClientRefreshTokenAge] = strconv.Itoa(client.RefreshTokenAge)
	return row
}

type RefreshToken struct {
	ClientId   string `json:"clientId"`
	Token      string `json:"string"`
	ExpireTime int64  `json:"expireTime"` // UNIX ms
}

func (r RefreshToken) StructToRow(refreshToken RefreshToken) []string {
	row := make([]string, 3)
	row[RefreshTokenClientId] = refreshToken.ClientId
	row[RefreshTokenToken] = refreshToken.Token
	row[RefreshTokenExpireTime] = strconv.FormatInt(refreshToken.ExpireTime, 10)
	return row
}

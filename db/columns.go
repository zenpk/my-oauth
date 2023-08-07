package db

import (
	"log"
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

type Table interface {
	ToRow() []string
	FromRow(row []string)
}

type User struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) ToRow() []string {
	row := make([]string, 3)
	row[UserUuid] = u.Uuid
	row[UserUsername] = u.Username
	row[UserPassword] = u.Password
	return row
}

func (u *User) FromRow(row []string) {
	u.Uuid = row[UserUuid]
	u.Username = row[UserUsername]
	u.Password = row[UserPassword]
}

type Client struct {
	Id              string `json:"id"`
	Secret          string `json:"secret"`
	Redirects       string `json:"redirects"`
	AccessTokenAge  int    `json:"accessTokenAge"`  // hour
	RefreshTokenAge int    `json:"RefreshTokenAge"` // hour
}

func (c *Client) ToRow() []string {
	row := make([]string, 4)
	row[ClientId] = c.Id
	row[ClientSecret] = c.Secret
	row[ClientRedirects] = c.Redirects
	row[ClientAccessTokenAge] = strconv.Itoa(c.AccessTokenAge)
	row[ClientRefreshTokenAge] = strconv.Itoa(c.RefreshTokenAge)
	return row
}

func (c *Client) FromRow(row []string) {
	var err error
	c.Id = row[ClientId]
	c.Secret = row[ClientSecret]
	c.Redirects = row[ClientRedirects]
	c.AccessTokenAge, err = strconv.Atoi(row[ClientAccessTokenAge])
	if err != nil {
		log.Fatalln(err)
	}
	c.RefreshTokenAge, err = strconv.Atoi(row[ClientAccessTokenAge])
	if err != nil {
		log.Fatalln(err)
	}
}

type RefreshToken struct {
	ClientId   string `json:"clientId"`
	Token      string `json:"string"`
	ExpireTime int64  `json:"expireTime"` // UNIX ms
}

func (r *RefreshToken) ToRow() []string {
	row := make([]string, 3)
	row[RefreshTokenClientId] = r.ClientId
	row[RefreshTokenToken] = r.Token
	row[RefreshTokenExpireTime] = strconv.FormatInt(r.ExpireTime, 10)
	return row
}

func (r *RefreshToken) FromRow(row []string) {
	var err error
	r.ClientId = row[RefreshTokenClientId]
	r.Token = row[RefreshTokenToken]
	r.ExpireTime, err = strconv.ParseInt(row[RefreshTokenExpireTime], 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
}

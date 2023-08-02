package db

const (
	UserUuid     = 0
	UserUsername = 1
	UserPassword = 2

	ClientId       = 0
	ClientSecret   = 1
	ClientRedirect = 2
	ClientOwner    = 3

	RefreshTokenToken      = 0
	RefreshTokenExpireTime = 1
)

type User struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Client struct {
	Id        string `json:"id"`
	Secret    string `json:"secret"`
	Redirects string `json:"redirects"`
	Owner     string `json:"owner"`
}

type RefreshToken struct {
	Token      string `json:"string"`
	ExpireTime string `json:"expireTime"` // UNIX ms
}

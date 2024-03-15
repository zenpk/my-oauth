package token

import (
	"github.com/cristalhq/jwt/v5"
	"github.com/zenpk/my-oauth/utils"
)

type Token struct {
	conf *utils.Configuration
}

type Claim struct {
	jwt.RegisteredClaims
	Uuid     string
	Username string
	ClientId string
}

func (t *Token) Init(conf *utils.Configuration){
	t.conf = conf
}

func mainn() {
	// create a Signer (HMAC in this example)
	key := []byte(`secret`)
	signer, err := jwt.NewSignerRS(jwt.RS256, key)

	// create claims (you can create your own, see: ExampleBuilder_withUserClaims)
	claims := &jwt.RegisteredClaims{
		Audience: []string{"admin"},
		ID:       "random-unique-string",
	}

	// create a Builder
	builder := jwt.NewBuilder(signer)

	// and build a Token
	token, err := builder.Build(claims)
	checkErr(err)

	// here is token as a string
	var _ string = token.String()
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

func GenAndInsertRefreshToken(dbInstance *db.Db, payload Payload, tokenAge time.Duration) (string, error) {
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

func GetAndCleanRefreshToken(dbInstance *db.Db, refreshToken string) (db.RefreshToken, error) {
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

package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/zenpk/my-oauth/util"
)

type Token struct {
	conf          *util.Configuration
	rsaPrivateKey *rsa.PrivateKey
}

type Claims struct {
	jwt.RegisteredClaims
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	ClientId string `json:"clientId"`
}

func (t *Token) Init(conf *util.Configuration) error {
	t.conf = conf
	// parse private key
	privateKeyBytes, err := os.ReadFile(t.conf.RsaPrivateKeyPath)
	if err != nil {
		return err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return errors.New("Token Init read RSA private key failed")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		// if it's not a PKCS#1 private key, try to parse it as a PKCS#8 key
		privateKeyInterface, err2 := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
		if err2 != nil {
			return err2
		}
		var ok bool
		privateKey, ok = privateKeyInterface.(*rsa.PrivateKey)
		if !ok {
			return errors.New("Token Init convert RSA private key failed")
		}
	}
	t.rsaPrivateKey = privateKey
	return nil
}

func (t *Token) GenerateJwt(claims *Claims) (string, error) {
	signer, err := jwt.NewSignerRS(jwt.RS256, t.rsaPrivateKey)
	if err != nil {
		return "", err
	}
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (t *Token) VerifyJwt(token *jwt.Token) (bool, error) {
	verifier, err := jwt.NewVerifierRS(jwt.RS256, &t.rsaPrivateKey.PublicKey)
	if err != nil {
		return false, err
	}
	// parse and verify a token
	tokenBytes := token.Bytes()
	parsedToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return false, err
	}
	newClaims := new(Claims)
	if err := json.Unmarshal(parsedToken.Claims(), newClaims); err != nil {
		return false, err
	}
	// don't want to bother implementing it
	// 	// at least match one audience
	// 	audiencePass := true
	// 	if requirement.Audience!=nil{
	// 	for _, aud := range requirement.Audience {
	// 		audiencePass = false
	// 		if newClaims.IsForAudience(aud) {
	// 			audiencePass = true
	// 			break
	// 		}
	// 	}
	// }
	// 	if !audiencePass {
	// 		log.Println("verify JWT error: audience not matched")
	// 		return false, nil
	// 	}
	if newClaims.Issuer != t.conf.JwtIssuer {
		log.Println("verify JWT error: issuer not matched")
		return false, nil
	}
	if newClaims.IsValidAt(time.Now()) {
		log.Println("verify JWT error: token expired")
		return false, nil
	}
	return true, nil
}

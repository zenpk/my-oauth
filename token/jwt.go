package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/zenpk/my-oauth/util"
)

type Claims struct {
	jwt.RegisteredClaims
	Uuid     string
	Username string
	ClientId string
}

func (t *Token) Init(conf *util.Configuration, logger *util.Logger) error {
	t.conf = conf
	t.logger = logger
	// parse private key
	pemData, err := os.ReadFile(t.conf.RsaPrivateKeyPath)
	if err != nil {
		return err
	}
	privateKeyBlock, _ := pem.Decode(pemData)
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

func (t *Token) GenJwt(claims *Claims) (string, error) {
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

func (t *Token) VerifyJwt(token string) (bool, error) {
	verifier, err := jwt.NewVerifierRS(jwt.RS256, &t.rsaPrivateKey.PublicKey)
	if err != nil {
		return false, err
	}
	// parse and verify a token
	parsedToken, err := jwt.Parse([]byte(token), verifier)
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
	// 		t.logger.Println("verify JWT error: audience not matched")
	// 		return false, nil
	// 	}
	if newClaims.Issuer != t.conf.JwtIssuer {
		t.logger.Println("verify JWT error: issuer not matched")
		return false, nil
	}
	if !newClaims.IsValidAt(time.Now()) {
		t.logger.Println("verify JWT error: token expired")
		return false, nil
	}
	return true, nil
}

type Jwk struct {
	Kty string `json:"kty"`
	E   string `json:"e"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
}

// Get converts an RSA public key in PEM format to a JWK
func (t *Token) GetJWK() (*Jwk, error) {
	// encode modulus and exponent to Base64 URL
	modulus := t.base64URLEncode(t.rsaPrivateKey.PublicKey.N)
	exponent := t.base64URLEncode(big.NewInt(int64(t.rsaPrivateKey.PublicKey.E)))
	return &Jwk{
		Kty: "RSA",
		E:   exponent,
		Use: "sig",
		Alg: "RS256",
		N:   modulus,
	}, nil
}

// base64URLEncode encodes a big integer (like RSA modulus or exponent) in the Base64 URL encoding
func (t *Token) base64URLEncode(value *big.Int) string {
	// the RawURLEncoding is used to avoid padding, which is typical for URL-encoded base64 variants
	return base64.RawURLEncoding.EncodeToString(value.Bytes())
}

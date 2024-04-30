package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/token"
	"github.com/zenpk/my-oauth/util"
)

type loginReq struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ClientId      string `json:"clientId"`
	CodeChallenge string `json:"codeChallenge"`
	Redirect      string `json:"redirect"`
}

type loginResp struct {
	commonResp
	AuthorizationCode string `json:"authorizationCode"`
}

func (h Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.Username == "" || req.Password == "" || req.ClientId == "" || req.CodeChallenge == "" || req.Redirect == "" {
		responseInputError(w)
		return
	}
	client, err := h.db.Clients.SelectByClientId(req.ClientId)
	if err != nil {
		responseError(w, err)
		return
	}
	if client == nil {
		responseErrMsg(w, "client id not found")
		return
	}
	user, err := h.db.Users.SelectByName(req.Username)
	if err != nil {
		responseError(w, err)
		return
	}
	if user == nil {
		responseErrMsg(w, "username doesn't exist")
		return
	}
	passwordMatch, err := util.BCryptHashCheck(user.Password, req.Password)
	if err != nil {
		responseError(w, err)
		return
	}
	if !passwordMatch {
		responseErrMsg(w, "incorrect password")
		return
	}
	redirects := strings.Split(client.Redirects, ",")
	redirectValid := false
	for _, redirect := range redirects {
		if strings.Trim(redirect, " ") == req.Redirect {
			redirectValid = true
			break
		}
	}
	if !redirectValid {
		responseErrMsg(w, "invalid redirect uri")
		return
	}
	authorizationCode, err := h.authInfo.GenAuthorizationCode(&util.AuthorizationInfo{
		ClientId:      client.Id,
		UserId:        user.Id,
		CodeChallenge: req.CodeChallenge,
	})
	if err != nil {
		responseError(w, err)
		return
	}
	// ok is ignored because it should always be true
	sw, _ := w.(*statusResponseWriter)
	sw.WriteUsername(user.Name)
	responseJson(sw, loginResp{
		commonResp:        genOkResponse(),
		AuthorizationCode: authorizationCode,
	})
}

type authorizeReq struct {
	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizationCode string `json:"authorizationCode"`
	CodeVerifier      string `json:"codeVerifier"`
}

type authorizeResp struct {
	commonResp
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h Handler) authorize(w http.ResponseWriter, r *http.Request) {
	var req authorizeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.ClientId == "" || req.ClientSecret == "" || req.AuthorizationCode == "" || req.CodeVerifier == "" {
		responseInputError(w)
		return
	}
	info, err := h.authInfo.VerifyAuthorizationCode(req.AuthorizationCode, req.CodeVerifier)
	if err != nil {
		responseErrMsg(w, err.Error())
		return
	}
	client, statusCode, err := h.checkClient(req.ClientId, req.ClientSecret)
	if err != nil {
		responseError(w, err, statusCode)
		return
	}
	if info.ClientId != client.Id {
		responseErrMsg(w, "client id not match")
		return
	}
	user, err := h.db.Users.SelectById(info.UserId)
	if err != nil {
		responseError(w, err)
		return
	}
	if user == nil {
		responseErrMsg(w, "user somehow doesn't exist")
		return
	}
	expireTime := time.Now().Add(time.Duration(client.AccessTokenAge) * time.Hour)
	claims := &token.Claims{
		Uuid:     user.Uuid,
		Username: user.Name,
		ClientId: client.ClientId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expireTime},
			Issuer:    h.conf.JwtIssuer,
		},
	}
	accessToken, err := h.tk.GenJwt(claims)
	if err != nil {
		responseError(w, err)
		return
	}
	refreshToken, err := h.sv.GenAndInsertRefreshToken(claims, client, user)
	if err != nil {
		responseError(w, err)
		return
	}
	sw, _ := w.(*statusResponseWriter)
	sw.WriteUsername(user.Name)
	responseJson(sw, authorizeResp{
		commonResp:   genOkResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

type refreshReq struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RefreshToken string `json:"refreshToken"`
}

type refreshResp struct {
	commonResp
	AccessToken string `json:"accessToken"`
}

func (h Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.ClientId == "" || req.ClientSecret == "" || req.RefreshToken == "" {
		responseInputError(w)
		return
	}
	client, statusCode, err := h.checkClient(req.ClientId, req.ClientSecret)
	if err != nil {
		responseError(w, err, statusCode)
		return
	}
	oldRefreshToken, err := h.sv.CleanAndGetRefreshToken(req.RefreshToken)
	if err != nil {
		responseError(w, err)
		return
	}
	if oldRefreshToken.ClientId != client.Id {
		responseErrMsg(w, "client id not match")
		return
	}
	// oldRefreshToken won't be nil here
	user, err := h.db.Users.SelectById(oldRefreshToken.UserId)
	if err != nil {
		responseError(w, err)
		return
	}
	if user == nil {
		responseErrMsg(w, "user somehow doesn't exist")
		return
	}
	expireTime := time.Now().Add(time.Duration(client.AccessTokenAge) * time.Hour)
	claims := &token.Claims{
		Uuid:     user.Uuid,
		Username: user.Name,
		ClientId: client.ClientId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expireTime},
			Issuer:    h.conf.JwtIssuer,
		},
	}
	accessToken, err := h.tk.GenJwt(claims)
	if err != nil {
		responseError(w, err)
		return
	}
	sw, _ := w.(*statusResponseWriter)
	sw.WriteUsername(user.Name)
	responseJson(sw, refreshResp{
		commonResp:  genOkResponse(),
		AccessToken: accessToken,
	})
}

type checkReq struct {
	AccessToken string `json:"accessToken"`
}

func (h Handler) verify(w http.ResponseWriter, r *http.Request) {
	var req checkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.AccessToken == "" {
		responseInputError(w)
		return
	}
	claims, ok, err := h.tk.ParseAndVerifyJwt(req.AccessToken)
	if err != nil {
		responseError(w, err)
		return
	}
	if !ok {
		responseErrMsg(w, "JWT parse failed")
		return
	}
	sw, _ := w.(*statusResponseWriter)
	sw.WriteUsername(claims.Username)
	responseOk(sw)
}

func (h Handler) checkClient(clientId, clientSecret string) (*dal.Client, int, error) {
	client, err := h.db.Clients.SelectByClientId(clientId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if client == nil {
		return nil, http.StatusOK, errors.New("client id doesn't exist")
	}
	secretMatch, err := util.BCryptHashCheck(client.Secret, clientSecret)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if !secretMatch {
		return nil, http.StatusOK, errors.New("incorrect client secret")
	}
	return client, http.StatusOK, nil
}

package handlers

import (
	"encoding/json"
	"errors"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/utils"
	"net/http"
	"strings"
	"time"
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
	client, err := h.Db.TableClient.Select(db.ClientId, req.ClientId)
	if err != nil {
		responseError(w, err)
		return
	}
	if client == nil {
		responseMsg(w, "client id not found")
		return
	}
	user, err := h.Db.TableUser.Select(db.UserUsername, req.Username)
	if err != nil {
		responseError(w, err)
		return
	}
	if user == nil {
		responseMsg(w, "username doesn't exist")
		return
	}
	passwordMatch, err := utils.BCryptHashCheck(user.(db.User).Password, req.Password)
	if err != nil {
		responseError(w, err)
		return
	}
	if !passwordMatch {
		responseMsg(w, "incorrect password")
		return
	}
	redirects := strings.Split(client.(db.Client).Redirects, ",")
	redirectValid := false
	for _, redirect := range redirects {
		if strings.Trim(redirect, " ") == req.Redirect {
			redirectValid = true
			break
		}
	}
	if !redirectValid {
		responseMsg(w, "invalid redirect uri")
		return
	}
	authorizationCode, err := utils.GenAuthorizationCode(utils.AuthorizationInfo{
		ClientId:      client.(db.Client).Id,
		Uuid:          user.(db.User).Uuid,
		Username:      user.(db.User).Username,
		CodeChallenge: req.CodeChallenge,
	})
	if err != nil {
		responseError(w, err)
		return
	}
	responseJson(w, loginResp{
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
	client, statusCode, err := h.checkClient(req.ClientId, req.ClientSecret)
	if err != nil {
		responseError(w, err, statusCode)
		return
	}
	info, err := utils.VerifyAuthorizationCode(req.AuthorizationCode, req.CodeVerifier)
	if err != nil {
		responseMsg(w, err.Error())
		return
	}
	if info.ClientId != client.Id {
		responseMsg(w, "client id not match")
		return
	}
	payload := utils.Payload{
		Uuid:     info.Uuid,
		Username: info.Username,
		ClientId: client.Id,
	}
	accessToken, err := utils.GenerateJwt(payload, time.Duration(client.RefreshTokenAge)*time.Hour)
	if err != nil {
		responseError(w, err)
		return
	}
	refreshToken, err := utils.GenAndInsertRefreshToken(h.Db, payload, time.Duration(client.RefreshTokenAge)*time.Hour)
	if err != nil {
		responseError(w, err)
		return
	}
	responseJson(w, authorizeResp{
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
	oldRefreshToken, err := utils.GetAndCleanRefreshToken(h.Db, req.RefreshToken)
	if err != nil {
		responseError(w, err)
		return
	}
	if oldRefreshToken.ClientId != client.Id {
		responseMsg(w, "client id not match")
		return
	}
	payload := utils.Payload{
		Uuid:     oldRefreshToken.Uuid,
		Username: oldRefreshToken.Username,
		ClientId: client.Id,
	}
	accessToken, err := utils.GenerateJwt(payload, time.Duration(client.RefreshTokenAge)*time.Hour)
	if err != nil {
		responseError(w, err)
		return
	}
	responseJson(w, refreshResp{
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
	if err := utils.VerifyJwt(req.AccessToken); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

func (h Handler) checkClient(clientId, clientSecret string) (db.Client, int, error) {
	client, err := h.Db.TableClient.Select(db.ClientId, clientId)
	if err != nil {
		return db.Client{}, http.StatusInternalServerError, err
	}
	if client == nil {
		return db.Client{}, http.StatusOK, errors.New("client id doesn't exist")
	}
	secretMatch, err := utils.BCryptHashCheck(client.(db.Client).Secret, clientSecret)
	if err != nil {
		return db.Client{}, http.StatusInternalServerError, err
	}
	if !secretMatch {
		return db.Client{}, http.StatusOK, errors.New("incorrect client secret")
	}
	return client.(db.Client), http.StatusOK, nil
}

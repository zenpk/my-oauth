package handlers

import (
	"encoding/json"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/utils"
	"net/http"
	"strings"
)

type loginReq struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ClientId      string `json:"clientId"`
	ClientSecret  string `json:"clientSecret"`
	CodeChallenge string `json:"codeChallenge"`
	RedirectUri   string `json:"redirectUri"`
}

type loginResp struct {
	commonResp
	AuthorizationCode string `json:"authorizationCode"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.Username == "" || req.Password == "" || req.ClientId == "" || req.ClientSecret == "" || req.CodeChallenge == "" || req.RedirectUri == "" {
		responseInputError(w)
		return
	}
	client, err := db.TableClient.Select(db.ClientId, req.ClientId)
	if err != nil {
		responseError(w, err)
		return
	}
	if client == nil {
		responseMsg(w, "client id doesn't exist")
		return
	}
	secretMatch, err := utils.BCryptHashCheck(client.(db.Client).Secret, req.ClientSecret)
	if err != nil {
		responseError(w, err)
		return
	}
	if !secretMatch {
		responseMsg(w, "incorrect client secret")
		return
	}
	user, err := db.TableUser.Select(db.UserUsername, req.Username)
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
		if strings.Trim(redirect, " ") == req.RedirectUri {
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
		CodeChallenge: req.CodeChallenge,
	})
	if err != nil {
		responseError(w, err)
		return
	}
	responseJson(w, loginResp{
		commonResp: commonResp{
			Ok:  true,
			Msg: "ok",
		},
		AuthorizationCode: authorizationCode,
	})
}

type authorizeReq struct {
	ClientId          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
	AuthorizationCode string `json:"authorizationCode"`
	CodeVerifier      string `json:"codeVerifier"`
}

type tokenResp struct {
	commonResp
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func authorize(w http.ResponseWriter, r *http.Request) {
	var req authorizeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.ClientId == "" || req.ClientSecret == "" || req.AuthorizationCode == "" || req.CodeVerifier == "" {
		responseInputError(w)
		return
	}
	client, err := db.TableClient.Select(db.ClientId, req.ClientId)
	if err != nil {
		responseError(w, err)
		return
	}
	if client == nil {
		responseMsg(w, "client id doesn't exist")
		return
	}
	secretMatch, err := utils.BCryptHashCheck(client.(db.Client).Secret, req.ClientSecret)
	if err != nil {
		responseError(w, err)
		return
	}
	if !secretMatch {
		responseMsg(w, "incorrect client secret")
		return
	}
}

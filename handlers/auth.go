package handlers

import (
	"encoding/json"
	"errors"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/utils"
	"net/http"
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
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if req.ClientId == "" || req.ClientSecret == "" || req.CodeChallenge == "" {
		responseInputError(w)
		return
	}
	client, err := db.TableClient.Select(db.ClientId, req.ClientId)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if client == nil {
		responseError(w, errors.New("client id doesn't exist"), http.StatusOK)
		return
	}
	user, err := db.TableUser.Select(db.UserUsername, req.Username)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		responseError(w, errors.New("username doesn't exist"), http.StatusOK)
		return
	}
	secretMatch, err := utils.BCryptHashCheck(client.(db.Client).Secret, req.ClientSecret)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if !secretMatch {
		responseError(w, errors.New("incorrect client secret"), http.StatusOK)
		return
	}
	passwordMatch, err := utils.BCryptHashCheck(user.(db.User).Password, req.Password)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if !passwordMatch {
		responseError(w, errors.New("incorrect password"), http.StatusOK)
		return
	}
	//
}

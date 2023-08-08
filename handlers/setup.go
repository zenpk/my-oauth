package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/utils"
	"net/http"
)

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("code") != utils.Conf.InvitationCode {
		responseError(w, errors.New("sorry, you need an invitation code or the code is wrong"), http.StatusOK)
		return
	}
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		responseInputError(w)
		return
	}
	res, err := db.UserCsv.Select(db.UserUsername, req.Username)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if res != nil {
		responseError(w, errors.New("user already exists"), http.StatusOK)
		return
	}
	passwordHash, err := utils.BCryptPassword(req.Password)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	user := db.User{
		Uuid:     uuid.New().String(),
		Username: req.Username,
		Password: passwordHash,
	}
	if err := db.UserCsv.Insert(user); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	responseOk(w)
}

type clientListResp struct {
	commonResp
	Clients []db.Client `json:"clients"`
}

func clientList(w http.ResponseWriter, r *http.Request) {
	clients, err := db.ClientCsv.All()
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	clientsConverted := make([]db.Client, 0)
	for _, client := range clients {
		converted := client.(db.Client)
		clientsConverted = append(clientsConverted, converted)
	}
	responseJson(w, clientListResp{
		commonResp: commonResp{
			Ok:  true,
			Msg: "ok",
		},
		Clients: clientsConverted,
	}, http.StatusOK)
}

type clientCreateReq struct {
	db.Client
	AdminPassword string `json:"adminPassword"`
}

func clientCreate(w http.ResponseWriter, r *http.Request) {
	var req clientCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if req.Id == "" || req.Secret == "" || req.Redirects == "" || req.AccessTokenAge <= 0 || req.RefreshTokenAge <= 0 || req.AdminPassword == "" {
		responseInputError(w)
		return
	}
	if req.RefreshTokenAge <= req.AccessTokenAge {
		responseError(w, errors.New("refresh_token age should be longer than access_token age"), http.StatusOK)
		return
	}
	passwordMatch, err := utils.BCryptHashCheck(utils.Conf.AdminPassword, req.AdminPassword)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if !passwordMatch {
		responseError(w, errors.New("incorrect admin password"), http.StatusOK)
		return
	}
	hashedSecret, err := utils.BCryptPassword(req.Secret)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	req.Secret = hashedSecret
	if err := db.ClientCsv.Insert(req.Client); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	responseOk(w)
}

type clientDeleteReq struct {
	ClientId      string `json:"clientId"`
	AdminPassword string `json:"adminPassword"`
}

func clientDelete(w http.ResponseWriter, r *http.Request) {
	var req clientDeleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if req.ClientId == "" || req.AdminPassword == "" {
		responseInputError(w)
		return
	}
	passwordMatch, err := utils.BCryptHashCheck(utils.Conf.AdminPassword, req.AdminPassword)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if !passwordMatch {
		responseError(w, errors.New("incorrect admin password"), http.StatusOK)
		return
	}
	if err := db.ClientCsv.Delete(db.ClientId, req.ClientId); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	responseOk(w)
}

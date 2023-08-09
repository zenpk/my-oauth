package handlers

import (
	"encoding/json"
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
		responseMsg(w, "sorry, you need an invitation code or the code is wrong")
		return
	}
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.Username == "" || req.Password == "" {
		responseInputError(w)
		return
	}
	res, err := db.TableUser.Select(db.UserUsername, req.Username)
	if err != nil {
		responseError(w, err)
		return
	}
	if res != nil {
		responseMsg(w, "user already exists")
		return
	}
	passwordHash, err := utils.BCryptPassword(req.Password)
	if err != nil {
		responseError(w, err)
		return
	}
	user := db.User{
		Uuid:     uuid.New().String(),
		Username: req.Username,
		Password: passwordHash,
	}
	if err := db.TableUser.Insert(user); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

type clientListResp struct {
	commonResp
	Clients []db.Client `json:"clients"`
}

func clientList(w http.ResponseWriter, r *http.Request) {
	clients, err := db.TableClient.All()
	if err != nil {
		responseError(w, err)
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
	})
}

type clientCreateReq struct {
	db.Client
	AdminPassword string `json:"adminPassword"`
}

func clientCreate(w http.ResponseWriter, r *http.Request) {
	var req clientCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.Id == "" || req.Secret == "" || req.Redirects == "" || req.AccessTokenAge <= 0 || req.RefreshTokenAge <= 0 || req.AdminPassword == "" {
		responseInputError(w)
		return
	}
	if req.RefreshTokenAge <= req.AccessTokenAge {
		responseMsg(w, "refresh_token age should be longer than access_token age")
		return
	}
	passwordMatch, err := utils.BCryptHashCheck(utils.Conf.AdminPassword, req.AdminPassword)
	if err != nil {
		responseError(w, err)
		return
	}
	if !passwordMatch {
		responseMsg(w, "incorrect admin password")
		return
	}
	hashedSecret, err := utils.BCryptPassword(req.Secret)
	if err != nil {
		responseError(w, err)
		return
	}
	req.Secret = hashedSecret
	if err := db.TableClient.Insert(req.Client); err != nil {
		responseError(w, err)
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
		responseInputError(w, err)
		return
	}
	if req.ClientId == "" || req.AdminPassword == "" {
		responseInputError(w)
		return
	}
	passwordMatch, err := utils.BCryptHashCheck(utils.Conf.AdminPassword, req.AdminPassword)
	if err != nil {
		responseError(w, err)
		return
	}
	if !passwordMatch {
		responseMsg(w, "incorrect admin password")
		return
	}
	if err := db.TableClient.Delete(db.ClientId, req.ClientId); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/utils"
	"net/http"
)

type registerReq struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	InvitationCode string `json:"invitationCode"`
}

func register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.InvitationCode != utils.Conf.InvitationCode {
		responseMsg(w, "sorry, you need an invitation code or the code is incorrect")
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

type clientWithoutSecret struct {
	Id              string `json:"id"`
	Redirects       string `json:"redirects"`
	AccessTokenAge  int    `json:"accessTokenAge"`
	RefreshTokenAge int    `json:"refreshTokenAge"`
}

type clientListResp struct {
	commonResp
	Clients []clientWithoutSecret `json:"clients"`
}

func clientList(w http.ResponseWriter, r *http.Request) {
	clients, err := db.TableClient.All()
	if err != nil {
		responseError(w, err)
		return
	}
	clientsConverted := make([]clientWithoutSecret, 0)
	for _, client := range clients {
		converted := clientWithoutSecret{
			Id:              client.(db.Client).Id,
			Redirects:       client.(db.Client).Redirects,
			AccessTokenAge:  client.(db.Client).AccessTokenAge,
			RefreshTokenAge: client.(db.Client).RefreshTokenAge,
		}
		clientsConverted = append(clientsConverted, converted)
	}
	responseJson(w, clientListResp{
		commonResp: genOkResponse(),
		Clients:    clientsConverted,
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
		responseMsg(w, "refresh token age should be longer than access token age")
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
	oldClient, err := db.TableClient.Select(db.ClientId, req.Id)
	if err != nil {
		responseError(w, err)
		return
	}
	if oldClient != nil {
		responseMsg(w, "client id already exists")
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
	Id            string `json:"id"`
	AdminPassword string `json:"adminPassword"`
}

func clientDelete(w http.ResponseWriter, r *http.Request) {
	var req clientDeleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.Id == "" || req.AdminPassword == "" {
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
	if err := db.TableClient.Delete(db.ClientId, req.Id); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

func publicKey(w http.ResponseWriter, r *http.Request) {
	responseJson(w, utils.Conf.JwtPublicKey)
}

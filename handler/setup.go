package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/zenpk/my-oauth/dal"
	"github.com/zenpk/my-oauth/util"
)

type registerReq struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	InvitationCode string `json:"invitationCode"`
}

func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.InvitationCode != h.conf.InvitationCode {
		responseErrMsg(w, "sorry, you need an invitation code or the code is incorrect")
		return
	}
	if req.Username == "" || req.Password == "" {
		responseInputError(w)
		return
	}
	if len(req.Password) < h.conf.PasswordMinLength {
		responseErrMsg(w, "the password should be at least "+strconv.Itoa(h.conf.PasswordMinLength)+" characters long")
		return
	}
	checkUser, err := h.db.Users.SelectByName(req.Username)
	if err != nil {
		responseError(w, err)
		return
	}
	if checkUser != nil {
		responseErrMsg(w, "user already exists")
		return
	}
	passwordHash, err := util.BCryptPassword(req.Password)
	if err != nil {
		responseError(w, err)
		return
	}
	user := &dal.User{
		Uuid:     uuid.New().String(),
		Name:     req.Username,
		Password: passwordHash,
	}
	if err := h.db.Users.Insert(user); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

type clientWithoutSecret struct {
	Id              int64  `json:"id"`
	ClientId        string `json:"clientId"`
	Redirects       string `json:"redirects"`
	AccessTokenAge  int64  `json:"accessTokenAge"`
	RefreshTokenAge int64  `json:"refreshTokenAge"`
}

type clientListResp struct {
	commonResp
	Clients []clientWithoutSecret `json:"clients"`
}

func (h Handler) clientList(w http.ResponseWriter, r *http.Request) {
	clients, err := h.db.Clients.SelectAll()
	if err != nil {
		responseError(w, err)
		return
	}
	clientsConverted := make([]clientWithoutSecret, 0)
	for _, client := range clients {
		converted := clientWithoutSecret{
			Id:              client.Id,
			ClientId:        client.ClientId,
			Redirects:       client.Redirects,
			AccessTokenAge:  client.AccessTokenAge,
			RefreshTokenAge: client.RefreshTokenAge,
		}
		clientsConverted = append(clientsConverted, converted)
	}
	responseJson(w, clientListResp{
		commonResp: genOkResponse(),
		Clients:    clientsConverted,
	})
}

type clientCreateReq struct {
	dal.Client
	AdminPassword string `json:"adminPassword"`
}

func (h Handler) clientCreate(w http.ResponseWriter, r *http.Request) {
	var req clientCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.ClientId == "" || req.Secret == "" || req.Redirects == "" || req.AccessTokenAge <= 0 || req.RefreshTokenAge <= 0 || req.AdminPassword == "" {
		responseInputError(w)
		return
	}
	if req.RefreshTokenAge <= req.AccessTokenAge {
		responseErrMsg(w, "refresh token age should be longer than access token age")
		return
	}
	passwordMatch, err := util.BCryptHashCheck(h.conf.AdminPassword, req.AdminPassword)
	if err != nil {
		responseError(w, err)
		return
	}
	if !passwordMatch {
		responseErrMsg(w, "incorrect admin password")
		return
	}
	oldClient, err := h.db.Clients.SelectByClientId(req.ClientId)
	if err != nil {
		responseError(w, err)
		return
	}
	if oldClient != nil {
		responseErrMsg(w, "client id already exists")
		return
	}
	hashedSecret, err := util.BCryptPassword(req.Secret)
	if err != nil {
		responseError(w, err)
		return
	}
	req.Secret = hashedSecret
	if err := h.db.Clients.Insert(&req.Client); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

type clientDeleteReq struct {
	Id            int64  `json:"id"`
	AdminPassword string `json:"adminPassword"`
}

func (h Handler) clientDelete(w http.ResponseWriter, r *http.Request) {
	var req clientDeleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseInputError(w, err)
		return
	}
	if req.Id <= 0 || req.AdminPassword == "" {
		responseInputError(w)
		return
	}
	passwordMatch, err := util.BCryptHashCheck(h.conf.AdminPassword, req.AdminPassword)
	if err != nil {
		responseError(w, err)
		return
	}
	if !passwordMatch {
		responseErrMsg(w, "incorrect admin password")
		return
	}
	if err := h.db.Clients.DeleteById(req.Id); err != nil {
		responseError(w, err)
		return
	}
	responseOk(w)
}

func (h Handler) publicKey(w http.ResponseWriter, r *http.Request) {
	jwk, err := h.tk.GetJWK()
	if err != nil {
		responseError(w, err)
		return
	}
	responseJson(w, jwk)
}

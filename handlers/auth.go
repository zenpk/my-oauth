package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/zenpk/my-oauth/db"
	"github.com/zenpk/my-oauth/utils"
	"net/http"
)

type

func login(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")
	clientSecret := r.URL.Query().Get("clientSecret")
	codeChallenge := r.URL.Query().Get("codeChallenge")
	state := r.URL.Query().Get("state")
	if clientId == "" || clientSecret == "" || codeChallenge == "" {
		responseError(w, errors.New("some parameters are missing"), http.StatusOK)
		return
	}
	var req registerReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	res, err := db.UserTable.Select(db.UserUsername, req.Username)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	if len(res) > 0 {
		responseError(w, errors.New("user already exists"), http.StatusBadRequest)
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
	if err := db.UserTable.Insert(user.StructToRow(user)); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	responseJson(w, commonResp{
		Ok:  true,
		Msg: "ok",
	}, http.StatusOK)
}

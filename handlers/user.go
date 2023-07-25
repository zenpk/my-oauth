package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	InvitationCode = "your_code"
)

func Register(w http.ResponseWriter, r *http.Request) {
	printLog("/register", r)
	var u user
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusBadRequest)
		return
	}
	maxId, _, err := findByUsername(u.Username)
	if err != ErrUserNotFound {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: "user already exists",
			},
			"",
		}, http.StatusOK)
		return
	}
	u.Id = maxId + 1
	if err := addUser(u); err != nil {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusOK)
		return
	}
	token, err := genBasicToken()
	if err != nil {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusBadRequest)
		return
	}
	response(w, tokenResp{
		resp: resp{
			Ok:  true,
			Msg: "ok",
		},
		Token: token,
	}, http.StatusOK)
}

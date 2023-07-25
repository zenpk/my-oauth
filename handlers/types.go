package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	printLog("/login", r)
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
	_, foundUser, err := findByUsername(u.Username)
	if err != nil {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusOK)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(u.Password)); errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: "wrong password",
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

func TokenGen(w http.ResponseWriter, r *http.Request) {
	printLog("/token-gen", r)
	var req tokenGenReq
	err := json.NewDecoder(r.Body).Decode(&req)
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
	_, err = parseBasicToken(req.Token)
	if err != nil {
		response(w, tokenResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusOK)
		return
	}
	token, err := genDataToken(req.AppId, req.Data, req.Age)
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

func TokenCheck(w http.ResponseWriter, r *http.Request) {
	printLog("/token-check", r)
	var req tokenCheckParseReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, resp{
			Ok:  false,
			Msg: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	claims, err := parseDataToken(req.Token)
	if err != nil {
		response(w, resp{
			Ok:  false,
			Msg: err.Error(),
		}, http.StatusOK)
		return
	}
	if req.AppId != claims.AppId {
		response(w, resp{
			Ok:  false,
			Msg: "wrong appId",
		}, http.StatusOK)
		return
	}
	response(w, resp{
		Ok:  true,
		Msg: "ok",
	}, http.StatusOK)
}

func TokenParse(w http.ResponseWriter, r *http.Request) {
	printLog("/token-parse", r)
	var req tokenCheckParseReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response(w, tokenParseResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusBadRequest)
		return
	}
	claims, err := parseDataToken(req.Token)
	if err != nil {
		response(w, tokenParseResp{
			resp{
				Ok:  false,
				Msg: err.Error(),
			},
			"",
		}, http.StatusOK)
		return
	}
	if req.AppId != claims.AppId {
		response(w, resp{
			Ok:  false,
			Msg: "wrong appId",
		}, http.StatusOK)
		return
	}
	response(w, tokenParseResp{
		resp{
			Ok:  true,
			Msg: "ok",
		},
		claims.Data,
	}, http.StatusOK)
}

func response(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
	log.Println(data)
}

func printLog(route string, r *http.Request) {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	log.Printf("%v %v\n", route, ipAddress)
}

package handlers

import (
	"encoding/json"
	"github.com/zenpk/my-oauth/db"
	"net/http"
)

type Handler struct {
	Db *db.Db
}

type commonResp struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

func responseJson(w http.ResponseWriter, data any, statusCodes ...int) {
	code := http.StatusOK
	w.Header().Set("Content-Type", "application/json")
	if len(statusCodes) > 0 {
		code = statusCodes[0]
	}
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func responseOk(w http.ResponseWriter) {
	responseJson(w, commonResp{
		Ok:  true,
		Msg: "ok",
	})
}

func responseMsg(w http.ResponseWriter, msg string) {
	responseJson(w, commonResp{
		Ok:  false,
		Msg: msg,
	})
}

func responseInputError(w http.ResponseWriter, errs ...error) {
	msg := "input error"
	if len(errs) > 0 {
		msg = errs[0].Error()
	}
	responseJson(w, commonResp{
		Ok:  false,
		Msg: msg,
	}, http.StatusBadRequest)
}

func responseError(w http.ResponseWriter, err error, statusCodes ...int) {
	code := http.StatusInternalServerError
	if len(statusCodes) > 0 {
		code = statusCodes[0]
	}
	responseJson(w, commonResp{
		Ok:  false,
		Msg: err.Error(),
	}, code)
}

func genOkResponse() commonResp {
	return commonResp{
		Ok:  true,
		Msg: "ok",
	}
}

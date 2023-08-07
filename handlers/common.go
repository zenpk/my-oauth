package handlers

import (
	"encoding/json"
	"net/http"
)

type commonResp struct {
	Ok  bool   `json:"ok"`
	Msg string `json:"msg"`
}

func responseJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func responseError(w http.ResponseWriter, err error, statusCode int) {
	responseJson(w, commonResp{
		Ok:  false,
		Msg: err.Error(),
	}, statusCode)
}

func responseInputError(w http.ResponseWriter) {
	responseJson(w, commonResp{
		Ok:  false,
		Msg: "input error",
	}, http.StatusOK)
}

func responseOk(w http.ResponseWriter) {
	responseJson(w, commonResp{
		Ok:  true,
		Msg: "ok",
	}, http.StatusOK)
}

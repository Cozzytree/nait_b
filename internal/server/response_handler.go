package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data   any `json:"data"`
	Status int `json:"status"`
}

// ResponseWithError writes an error response with the provided message and status code.
func ResponseWithError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, msg, status)
}

func ResponseWithJson(w http.ResponseWriter, res Response) {
	w.WriteHeader(res.Status)
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(&res)
	if err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError)
	}
	w.Write(data)
}

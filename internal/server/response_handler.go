package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data   any `json:"data"`
	Status int `json:"status"`
}

// ResponseWithError writes an error response with the provided message and status code.
func ResponseWithError(w http.ResponseWriter, msg string, status int) {
	// Define a struct to hold the error message
	type res struct {
		Error string `json:"error"`
	}

	// Create the response object
	r := res{
		Error: fmt.Sprintf("error : %v", msg),
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(&r)
	if err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError)
	}
	w.Write(data)
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

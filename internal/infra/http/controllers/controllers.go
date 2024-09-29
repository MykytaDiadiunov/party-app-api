package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ctxKey struct {
	name string
}

var (
	UserKey    = ctxKey{"user"}
	SessionKey = ctxKey{"session"}
)

func Ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Success(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(err)
	}
}

func Created(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Print(err)
	}
}

func BadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	encodeErrorData(w, err)
}

func InternalServerError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	encodeErrorData(w, err)
}

func NotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	encodeErrorData(w, err)
}

func Unauthorized(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	encodeErrorData(w, err)
}

func encodeErrorData(w http.ResponseWriter, err error) {
	e := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if e != nil {
		log.Print(e)
	}
}

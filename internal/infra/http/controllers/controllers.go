package controllers

import (
	"context"
	"encoding/json"
	"go-rest-api/internal/domain"
	"log"
	"net/http"
)

type ctxKey struct {
	name string
}

type CtxStrKey string

var (
	UserKey    = ctxKey{"user"}
	SessionKey = ctxKey{"session"}
	PartyKey   = ctxKey{"party"}
)

func GetPathValueInCtx[T any](ctx context.Context, value T) context.Context {
	key := ResolveCtxKeyFromPathType(new(T))
	return context.WithValue(ctx, key, value)
}

func GetPathValueFromCtx[T any](ctx context.Context) T {
	key := ResolveCtxKeyFromPathType(new(T))
	return ctx.Value(key).(T)
}

func (k CtxStrKey) UrlParam() string {
	return string(k) + "Id"
}

func ResolveCtxKeyFromPathType(value any) CtxStrKey {
	switch value.(type) {
	case *domain.Party:
		return CtxStrKey(PartyKey.name)
	default:
		panic("unk type in resolveCtxKeyFromPathType (controller)")
	}
}

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

func NoContent(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	encodeErrorData(w, err)
}

func Unauthorized(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	encodeErrorData(w, err)
}

func Forbidden(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	encodeErrorData(w, err)
}

func encodeErrorData(w http.ResponseWriter, err error) {
	e := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if e != nil {
		log.Print(e)
	}
}

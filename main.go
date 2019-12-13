package main

import (
	"encoding/json"
	"errors"
	jwt "github.com/gbrlsnchs/jwt/v3"
	"github.com/mistyfiky/agh-sr-hades/model"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/ping",
		errorMiddleware(
			methodMiddleware("GET",
				corsMiddleware(
					pingHandler()))))
	http.HandleFunc("/authenticate",
		errorMiddleware(
			methodMiddleware("POST",
				corsMiddleware(
					authenticateHandler()))))
	if err := http.ListenAndServe(":80", nil); nil != err {
		panic(err)
	}
}

func respond(writer http.ResponseWriter, statusCode int, response interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if nil == response {
		return
	}
	body, err := json.Marshal(response)
	if nil != err {
		panic(err)
	}
	if _, err := writer.Write(body); nil != err {
		panic(err)
	}
}

func errorMiddleware(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if result := recover(); nil != result {
				err, ok := result.(error)
				if !ok {
					err = errors.New("undefined error")
				}
				response := model.SimpleResponse{
					Meta: model.Meta{
						Success: false,
						Message: err.Error(),
					},
				}
				respond(writer, http.StatusInternalServerError, response)
			}
		}()
		next.ServeHTTP(writer, request)
	}
}

func methodMiddleware(method string, next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if method != request.Method {
			response := model.SimpleResponse{
				Meta: model.Meta{
					Success: false,
					Message: http.StatusText(http.StatusMethodNotAllowed),
				},
			}
			respond(writer, http.StatusMethodNotAllowed, response)
			return
		}
		next.ServeHTTP(writer, request)
	}
}

func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "content-type")
		if "OPTIONS" == request.Method {
			respond(writer, http.StatusOK, nil)
			return
		}
		next.ServeHTTP(writer, request)
	}
}

func pingHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response := model.SimpleResponse{
			Meta: model.Meta{
				Success: true,
				Message: "pong!",
			},
		}
		respond(writer, http.StatusOK, response)
	}
}

func authenticateHandler() http.HandlerFunc {
	alg := jwt.NewHS256([]byte("secret"))
	issuer := "hades"
	return func(writer http.ResponseWriter, request *http.Request) {
		var auth model.Auth
		if err := json.NewDecoder(request.Body).Decode(&auth); nil != err {
			response := model.SimpleResponse{
				Meta: model.Meta{
					Success: false,
					Message: err.Error(),
				},
			}
			respond(writer, http.StatusBadRequest, response)
			return
		}
		payload := jwt.Payload{
			Subject:  auth.Username,
			Issuer:   issuer,
			IssuedAt: jwt.NumericDate(time.Now()),
		}
		token, err := jwt.Sign(payload, alg)
		if nil != err {
			panic(err)
		}
		response := model.TokenResponse{
			Meta: model.Meta{
				Success: true,
				Message: "success",
			},
			Data: model.Token{
				Token: string(token[:]),
			},
		}
		respond(writer, http.StatusOK, response)
	}
}

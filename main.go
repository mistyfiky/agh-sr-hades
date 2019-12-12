package main

import (
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
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func respond(writer http.ResponseWriter, statusCode int, body model.Response) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if statusCode > 299 && nil == body {
		body = model.NewResponseError(http.StatusText(statusCode))
	}
	if nil != body {
		if _, err := writer.Write(body.ToJson()); nil != err {
			panic(err)
		}
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
				respond(writer, http.StatusInternalServerError, model.NewResponseError(err.Error()))
			}
		}()
		next.ServeHTTP(writer, request)
	}
}

func methodMiddleware(method string, next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if method != request.Method {
			respond(writer, http.StatusMethodNotAllowed, nil)
		} else {
			next.ServeHTTP(writer, request)
		}
	}
}

func corsMiddleware(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "content-type")
		if "OPTIONS" == request.Method {
			respond(writer, http.StatusOK, nil)
		} else {
			next.ServeHTTP(writer, request)
		}
	}
}

func pingHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body := model.NewSimpleResponse("pong!").ToJson()
		if _, err := writer.Write(body); nil != err {
			panic(err)
		}
	}
}

func authenticateHandler() http.HandlerFunc {
	var alg = jwt.NewHS256([]byte("secret"))
	return func(writer http.ResponseWriter, request *http.Request) {
		payload := jwt.Payload{
			Subject:  "someone",
			Issuer:   "hades",
			IssuedAt: jwt.NumericDate(time.Now()),
		}
		token, err := jwt.Sign(payload, alg)
		if nil != err {
			panic(err)
		}
		body := model.NewTokenResponse(token).ToJson()
		if _, err := writer.Write(body); nil != err {
			panic(err)
		}
	}
}

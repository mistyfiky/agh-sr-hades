package main

import (
	jwt "github.com/gbrlsnchs/jwt/v3"
	"github.com/mistyfiky/agh-sr-hades/model"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/ping",
		contentTypeMiddleware(
			corsMiddleware(
				pingHandler())))
	http.HandleFunc("/authenticate",
		contentTypeMiddleware(
			corsMiddleware(
				authenticateHandler())))
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func contentTypeMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func corsMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := model.NewSimpleResponse("pong!").ToJson()
		if _, err := w.Write(body); err != nil {
			panic(err)
		}
	}
}

func authenticateHandler() http.HandlerFunc {
	var alg = jwt.NewHS256([]byte("secret"))
	return func(w http.ResponseWriter, r *http.Request) {
		payload := jwt.Payload{
			Subject:  "someone",
			Issuer:   "hades",
			IssuedAt: jwt.NumericDate(time.Now()),
		}
		token, err := jwt.Sign(payload, alg)
		if err != nil {
			panic(err)
		}
		body := model.NewTokenResponse(token).ToJson()
		if _, err := w.Write(body); err != nil {
			panic(err)
		}
	}
}

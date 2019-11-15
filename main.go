package main

import (
    "encoding/json"
    "fmt"
    jwt "github.com/gbrlsnchs/jwt/v3"
    "net/http"
    "time"
)

type Meta struct {
    Message string `json:"message"`
}

type Response struct {
    Meta Meta `json:"meta"`
}

var hs = jwt.NewHS256([]byte("secret"))

func main() {
    http.HandleFunc("/authenticate", corsHandler(authenticateHandler()))

    http.ListenAndServe(":80", nil)
}

func corsHandler(h http.Handler) http.HandlerFunc {
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

func authenticateHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Conten-Type", "application/json")

        pl := jwt.Payload{
            Subject:  "someone",
            Issuer:   "hades",
            IssuedAt: jwt.NumericDate(time.Now()),
        }

        token, err := jwt.Sign(pl, hs)
        if err != nil {
            fmt.Print(fmt.Errorf("%v", err))
        }

        res := &Response{Meta: Meta{Message: string(token[:])}}
        j, err := json.Marshal(res)
        if err != nil {
            panic(err)
        }
        w.Write(j)
    }
}

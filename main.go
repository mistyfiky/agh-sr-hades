package main

import (
	"fmt"
	jwt "github.com/gbrlsnchs/jwt/v3"
	"net/http"
	"time"
)

var hs = jwt.NewHS256([]byte("secret"))

func main() {

	http.HandleFunc("/authenticate", authenticate)

	http.ListenAndServe(":80", nil)
}

func authenticate(writer http.ResponseWriter, request *http.Request) {

	pl := jwt.Payload{
        Subject:  "someone",
        Issuer:   "hades",
        IssuedAt: jwt.NumericDate(time.Now()),
    }

	token, err := jwt.Sign(pl, hs)
	if err != nil {
		fmt.Print(fmt.Errorf("%v", err))
	}

	fmt.Fprintf(writer, "%v", string(token[:]))
}

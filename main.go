package main

import (
	"context"
	"encoding/json"
	"errors"
	jwt "github.com/gbrlsnchs/jwt/v3"
	"github.com/mistyfiky/agh-sr-hades/model"
	"github.com/mistyfiky/agh-sr-hades/repository"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var addr string
var alg *jwt.HMACSHA

func init() {
	addr = ":" + os.Getenv("APP_PORT")
	alg = jwt.NewHS256([]byte(os.Getenv("JWT_KEY")))
}

func main() {
	http.HandleFunc("/ping",
		errorMiddleware(
			corsMiddleware(
				methodMiddleware("GET",
					pingHandler()))))
	http.HandleFunc("/authenticate",
		errorMiddleware(
			corsMiddleware(
				methodMiddleware("POST",
					authenticateHandler()))))
	http.HandleFunc("/users/",
		errorMiddleware(
			corsMiddleware(
				jwtMiddleware(
					usersHandler()))))
	log.Println("starting server on " + addr)
	if err := http.ListenAndServe(addr, nil); nil != err {
		log.Fatal(err.Error())
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

func getPayload(request *http.Request) jwt.Payload {
	payload, ok := request.Context().Value("jwt").(jwt.Payload)
	if !ok {
		panic("invalid payload")
	}
	return payload
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
				log.Println(err.Error())
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

func jwtMiddleware(next http.Handler) http.HandlerFunc {
	respondUnauthorized := func(writer http.ResponseWriter) {
		response := model.SimpleResponse{
			Meta: model.Meta{
				Success: false,
				Message: http.StatusText(http.StatusUnauthorized),
			},
		}
		respond(writer, http.StatusUnauthorized, response)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		if len(header) < 7 {
			respondUnauthorized(writer)
			return
		}
		token := header[7:]
		payload := jwt.Payload{}
		if _, err := jwt.Verify([]byte(token), alg, &payload); nil != err {
			respondUnauthorized(writer)
			return
		}
		ctx := context.WithValue(request.Context(), "jwt", payload)
		next.ServeHTTP(writer, request.WithContext(ctx))
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
		user := repository.FindUserByUsername(auth.Username)
		if nil == user || !user.IsAuthenticatedBy(auth.Password) {
			response := model.SimpleResponse{
				Meta: model.Meta{
					Success: false,
					Message: "Invalid username or password",
				},
			}
			respond(writer, http.StatusUnauthorized, response)
			return
		}
		payload := jwt.Payload{
			Subject:  user.GetUsername(),
			Issuer:   "hades",
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

func usersHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		segments := strings.Split(request.URL.Path, "/")
		switch  {
		case len(segments) == 3 && segments[1] == "users" && "GET" == request.Method:
			usersGet(writer, request)
		case len(segments) == 5 && segments[1] == "users" && segments[3] == "movies" && "PUT" == request.Method:
			usersMoviesPut(writer, request)
		default:
			http.NotFound(writer, request)
		}
	}
}

func usersGet(writer http.ResponseWriter, request *http.Request) {
	segments := strings.Split(request.URL.Path, "/")
	if len(segments) < 3 {
		panic("invalid url path")
	}
	username := segments[2]
	user := repository.FindUserByUsername(username)
	if nil == user {
		http.NotFound(writer, request)
		return
	}
	userMovies := repository.FindUserMoviesByUsername(username)
	movies := []string{}
	for _, userMovie := range userMovies {
		movies = append(movies, userMovie.GetMovieId())
	}
	response := model.UserResponse{
		Meta: model.Meta{
			Success: true,
			Message: "success",
		},
		Data: model.User{
			Username: user.GetUsername(),
			Movies: movies,
		},
	}
	respond(writer, http.StatusOK, response)
}

func usersMoviesPut(writer http.ResponseWriter, request *http.Request) {
	segments := strings.Split(request.URL.Path, "/")
	if len(segments) < 5 {
		panic("invalid url path")
	}
	username := segments[2]
	user := repository.FindUserByUsername(username)
	if nil == user {
		http.NotFound(writer, request)
		return
	}
	movieId := segments[4]
	userMovie := repository.FindUserMovieByUsernameAndMovieId(username, movieId)
	if nil == userMovie {
		userMovie = repository.NewUserMovie(username, movieId)
		repository.SaveUserMovie(userMovie)
	}
	respond(writer, http.StatusNoContent, nil)
}

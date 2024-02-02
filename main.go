package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AlpsMonaco/go-http-middleware/middleware"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func httpHandleHello(w http.ResponseWriter, r *http.Request) {
	logger.Println("Hello")
	w.Write([]byte("Hello"))
}

func httpHandleWorld(w http.ResponseWriter, r *http.Request) {
	logger.Println("World")
	w.Write([]byte("World"))
}

func LogMiddleware1(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("log middleware 1 in")
		next(w, r)
		logger.Println("log middleware 1 out")
	}
}

func LogMiddleware2(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("log middleware 2 in")
		next(w, r)
		logger.Println("log middleware 2 out")
	}
}

func main() {
	http := middleware.DefaultHTTPBuilder()
	http.With(LogMiddleware1, LogMiddleware2)
	http.HandleFunc("/hello", httpHandleHello)
	http.HandleFunc("/world", httpHandleWorld)
	http.ListenAndServe(":33333", nil)
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AlpsMonaco/go-http-middleware/middleware"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func HttpHandleHello(w http.ResponseWriter, r *http.Request) {
	logger.Println("Hello")
	w.Write([]byte("Hello"))
}

func HttpHandleWorld(w http.ResponseWriter, r *http.Request) {
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
	b := middleware.NewBuilder(&http.ServeMux{})
	b.With(LogMiddleware1, LogMiddleware2)
	b.HandleFunc("/hello", HttpHandleHello)
	b.ListenAndServe(":33333", nil)

	http := middleware.DefaultHTTPBuilder()
	http.With(LogMiddleware1, LogMiddleware2)
	http.HandleFunc("/hello", HttpHandleHello)
	http.HandleFunc("/world", HttpHandleWorld)
	http.ListenAndServe(":33333", nil)
}

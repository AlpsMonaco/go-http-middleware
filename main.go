package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

func HttpHandlePanic(w http.ResponseWriter, r *http.Request) {
	panic("panic when access /panic")
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

func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.WriteHeader(500)
				logger.Println(err)
			}
		}()
		next(w, r)
	}
}

func TimerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		next(w, r)
		logger.Printf("time elapsed:%dms\n", time.Now().Sub(st).Milliseconds())
	}
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s\n", r.RemoteAddr, r.RequestURI)
		next(w, r)
	}
}

func main() {
	http := middleware.DefaultHTTPBuilder()
	http.With(TimerMiddleware, RecoverMiddleware, LogMiddleware, LogMiddleware1, LogMiddleware2)
	http.HandleFunc("/panic", HttpHandlePanic)
	http.HandleFunc("/hello", HttpHandleHello)
	http.HandleFunc("/world", HttpHandleWorld)
	http.ListenAndServe(":33333", nil)
}

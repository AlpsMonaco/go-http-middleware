package middleware

import (
	"net/http"
)

type Middleware func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Builder struct {
	mux         Mux
	middlewares []Middleware
}

func (b *Builder) With(m ...Middleware) *Builder {
	b.middlewares = append(b.middlewares, m...)
	return b
}

func (b *Builder) Handle(pattern string, handler http.Handler) *Builder {
	b.mux.Handle(pattern, compileHandlerWithMiddleware(b.middlewares, handler.ServeHTTP))
	return b
}

func (b *Builder) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Builder {
	b.mux.HandleFunc(pattern, compileHandlerWithMiddleware(b.middlewares, handler))
	return b
}

func compileHandlerWithMiddleware(middlewares []Middleware, f http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

var defaultBuilder = Builder{http.DefaultServeMux, []Middleware{}}

func Add(m ...Middleware) {
	defaultBuilder.With(m...)
}

func Handle(pattern string, handler http.Handler) {
	defaultBuilder.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	defaultBuilder.HandleFunc(pattern, handler)
}

// var middlewares []Middleware

// func AddMiddleware(m Middleware) {
// 	middlewares = append(middlewares, m)
// }

// func With(m ...Middleware) {
// 	middlewares = append(middlewares, m...)
// }

// func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
// 	http.HandleFunc(pattern, compileHandlerWithMiddleware(middlewares, handler))
// }

// type httpHandlerWrapper struct {
// 	h http.Handler
// }

// func (hhw *httpHandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	compileHandlerWithMiddleware(middlewares, func(w http.ResponseWriter, r *http.Request) {
// 		hhw.h.ServeHTTP(w, r)
// 	})
// }

// func Handle(pattern string, handler http.Handler) {
// 	http.Handle(pattern, &httpHandlerWrapper{handler})
// }

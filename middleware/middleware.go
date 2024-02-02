package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

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

func (b *Builder) Middleware(m ...Middleware) *Builder {
	return b.With(m...)
}

func (b *Builder) Handle(pattern string, handler http.Handler) *Builder {
	b.mux.Handle(pattern, compileHandlerWithMiddleware(b.middlewares, handler.ServeHTTP))
	return b
}

func (b *Builder) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) *Builder {
	b.mux.HandleFunc(pattern, compileHandlerWithMiddleware(b.middlewares, handler))
	return b
}

func (b *Builder) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

func (b *Builder) ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}

func compileHandlerWithMiddleware(middlewares []Middleware, f http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

var defaultBuilder = Builder{http.DefaultServeMux, []Middleware{}}

func DefaultHTTPBuilder() *Builder {
	return &defaultBuilder
}

func Add(m ...Middleware) {
	defaultBuilder.With(m...)
}

func Handle(pattern string, handler http.Handler) {
	defaultBuilder.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	defaultBuilder.HandleFunc(pattern, handler)
}

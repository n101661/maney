package iris

import (
	"github.com/kataras/iris/v12"
)

type Server struct {
	app *iris.Application
}

func NewServer() *Server {
	s := &Server{
		app: iris.Default(),
	}

	s.registerRoutes()

	return s
}

func (s *Server) ListenAndServe(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	return s.app.Run(iris.TLS(addr, certFile, keyFile))
}

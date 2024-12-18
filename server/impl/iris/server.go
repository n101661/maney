package iris

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"

	"github.com/n101661/maney/database"
	authV2 "github.com/n101661/maney/pkg/services/auth"
	"github.com/n101661/maney/server/impl/iris/auth"
)

type Config struct {
	// SecretKey is for JWS.
	SecretKey         []byte
	PasswordSaltRound int
}

type Server struct {
	app         *iris.Application
	authService authV2.Service
	auth        *auth.Authentication

	db database.DB
}

func NewServer(cfg Config, authService authV2.Service) *Server {
	s := &Server{
		app:         newIrisApplication(),
		authService: authService,
		auth: auth.NewAuthentication(
			cfg.SecretKey,
			auth.WithPasswordSaltRound(cfg.PasswordSaltRound),
		),
		db: nil, // TODO
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

func newIrisApplication() *iris.Application {
	app := iris.Default()
	app.Validator = validator.New()

	cfg := iris.DefaultConfiguration()
	cfg.DisablePathCorrection = true
	cfg.DisablePathCorrectionRedirection = true
	app.Configure(iris.WithConfiguration(cfg))
	return app
}

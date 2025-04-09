package iris

import (
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/n101661/maney/database"
	"github.com/n101661/maney/server/impl/iris/auth"
	"github.com/n101661/maney/server/impl/iris/config"
	"github.com/n101661/maney/server/users"
)

type Config struct {
	LogLevel    config.LogLevel `toml:"log-level" comment:"Log level. It can be one of the following: debug, info, warn, error, fatal, disable. The default is info."`
	CorsOrigins []string        `toml:"cors-origin"`

	// SecretKey is for JWS.
	SecretKey         []byte `toml:"-"`
	PasswordSaltRound int    `toml:"-"`
}

type Server struct {
	app  *iris.Application
	auth *auth.Authentication

	userController *users.IrisController

	db database.DB
}

func NewServer(cfg *Config, userController *users.IrisController) *Server {
	s := &Server{
		app: newIrisApplication(cfg),
		auth: auth.NewAuthentication(
			cfg.SecretKey,
			auth.WithPasswordSaltRound(cfg.PasswordSaltRound),
		),
		userController: userController,
		db:             nil, // TODO
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

func newIrisApplication(config *Config) *iris.Application {
	app := iris.New()
	app.Validator = validator.New()

	app.Logger().SetLevel(string(config.LogLevel))
	app.Logger().Debugf("Log level set to %s", config.LogLevel)

	app.UseRouter(requestid.New())
	app.Logger().Debug("Using <UUID4> to identify requests")

	app.Use(recover.New())

	allowedOrigins := slices.Clone(config.CorsOrigins)
	if len(allowedOrigins) == 0 {
		allowedOrigins = append(allowedOrigins, "*")
	}
	app.Use(cors.New().
		ExtractOriginFunc(cors.DefaultOriginExtractor).
		AllowOrigins(allowedOrigins...).
		Handler())

	return app
}

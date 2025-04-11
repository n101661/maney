package iris

import (
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/middleware/requestid"

	"github.com/n101661/maney/database"
	"github.com/n101661/maney/server/impl/iris/auth"
	"github.com/n101661/maney/server/impl/iris/config"
	"github.com/n101661/maney/server/middleware/errors"
	"github.com/n101661/maney/server/middleware/logger"
	"github.com/n101661/maney/server/middleware/recover"
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

	app.Logger().TimeFormat = "2006-01-02 15:04:05.999"
	app.Logger().SetLevel(string(config.LogLevel))
	app.Logger().Debugf("Log level set to %s", config.LogLevel)

	allowedOrigins := parseAllowedOrigins(config.CorsOrigins)

	app.UseRouter(
		requestid.New(),
		logger.New(
			logger.WithRequestIDFunc(func(ctx iris.Context) string {
				id, _ := ctx.GetID().(string)
				return id
			}),
			logger.WithExcludeRequest(func(ctx iris.Context) bool {
				_, ok := excludedRequestPath[ctx.Path()]
				return ok
			}),
		),
		recover.New(),
		cors.New().
			ExtractOriginFunc(cors.DefaultOriginExtractor).
			AllowOrigins(allowedOrigins...).
			Handler(),
	)
	app.Logger().Debug("Using <UUID4> to identify requests")
	app.Logger().Debug("Allowed CORS origins: ", allowedOrigins)

	app.UseError(errors.HideInternalErrorHandler)

	return app
}

func parseAllowedOrigins(origins []string) []string {
	allowedOrigins := slices.Clone(origins)
	if len(allowedOrigins) == 0 {
		allowedOrigins = append(allowedOrigins, "*")
	}
	return allowedOrigins
}

var excludedRequestPath = map[string]struct{}{
	"/auth/refresh": {},
	"/login":        {},
	"/auth/logout":  {},
	"/sign-up":      {},
}

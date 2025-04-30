package iris

func (s *Server) registerRoutes() {

	s.app.Post("/auth/refresh", s.controllers.User.RefreshAccessToken)
	s.app.Post("/login", s.controllers.User.Login)
	s.app.Post("/auth/logout", s.controllers.User.Logout)
	s.app.Post("/sign-up", s.controllers.User.SignUp)

	user := s.app.Party("/", s.controllers.User.ValidateAccessToken)

	{ // user's config
		user.Put("/config", s.controllers.User.UpdateUserConfig)
		user.Get("/config", s.controllers.User.GetUserConfig)
	}
	{ // user's accounts
		user.Post("/accounts", s.controllers.Account.Create)
		user.Get("/accounts", s.controllers.Account.List)
		user.Put("/accounts/{accountId}", s.controllers.Account.Update)
		user.Delete("/accounts/{accountId}", s.controllers.Account.Delete)
	}
	{ // user's categories
		user.Post("/categories", s.controllers.Category.Create)
		user.Get("/categories", s.controllers.Category.List)
		user.Put("/categories/{categoryId}", s.controllers.Category.Update)
		user.Delete("/categories/{categoryId}", s.controllers.Category.Delete)
	}
	{ // user's shops
		user.Post("/shops", s.controllers.Shop.Create)
		user.Get("/shops", s.controllers.Shop.List)
		user.Put("/shops/{shopId}", s.controllers.Shop.Update)
		user.Delete("/shops/{shopId}", s.controllers.Shop.Delete)
	}
	{ // user's fees
		user.Post("/fees")
		user.Get("/fees")
		user.Put("/fees/{feeId}")
		user.Delete("/fees/{feeId}")
	}
	{ // user's daily items
		user.Post("/daily-items")
		user.Get("/daily-items")
		user.Put("/daily-items/{dailyItemId}")
		user.Delete("/daily-items/{dailyItemId}")
	}
}

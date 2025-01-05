package iris

func (s *Server) registerRoutes() {

	s.app.Post("/auth/refresh")
	s.app.Post("/login", s.Login)
	s.app.Post("/auth/logout", s.Logout)
	s.app.Post("/sign-up", s.SignUp)

	user := s.app.Party("/", s.ValidateAccessToken)

	{ // user's config
		user.Put("/config", s.UpdateConfig)
		user.Get("/config", s.GetConfig)
	}
	{ // user's accounts
		user.Post("/accounts", s.CreateAccount)
		user.Get("/accounts", s.ListAccounts)
		user.Put("/accounts/{accountId}", s.UpdateAccount)
		user.Delete("/accounts/{accountId}", s.DeleteAccount)
	}
	{ // user's categories
		user.Post("/categories")
		user.Get("/categories")
		user.Put("/categories/{categoryId}")
		user.Delete("/categories/{categoryId}")
	}
	{ // user's shops
		user.Post("/shops")
		user.Get("/shops")
		user.Put("/shops/{shopId}")
		user.Delete("/shops/{shopId}")
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

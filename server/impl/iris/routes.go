package iris

func (s *Server) registerRoutes() {

	s.app.Post("/log-in", s.LogIn)
	s.app.Post("/log-out", s.LogOut)
	s.app.Post("/sign-up", s.SignUp)

	user := s.app.Party("/users", s.auth.ValidateToken)

	{ // user's config
		user.Put("/config", s.UpdateConfig)
		user.Get("/config", s.GetConfig)
	}
	{ // user's accounts
		user.Post("/accounts", s.CreateAccount)
		user.Get("/accounts", s.ListAccounts)
		user.Put("/accounts/{oid}", s.UpdateAccount)
		user.Delete("/accounts/{oid}", s.DeleteAccount)
	}
	{ // user's categories
		user.Post("/categories")
		user.Get("/categories")
		user.Put("/categories/{oid}")
		user.Delete("/categories/{oid}")
	}
	{ // user's shops
		user.Post("/shops")
		user.Get("/shops")
		user.Put("/shops/{oid}")
		user.Delete("/shops/{oid}")
	}
	{ // user's fees
		user.Post("/fees")
		user.Get("/fees")
		user.Put("/fees/{oid}")
		user.Delete("/fees/{oid}")
	}
	{ // user's daily items
		user.Post("/daily-items")
		user.Get("/daily-items")
		user.Put("/daily-items/{oid}")
		user.Delete("/daily-items/{oid}")
		user.Get("/daily-items/export")
		user.Post("/daily-items/import")
	}
	{ // user's repeating items
		user.Post("/repeating-items")
		user.Get("/repeating-items")
		user.Put("/repeating-items/{oid}")
		user.Delete("/repeating-items/{oid}")
	}
}

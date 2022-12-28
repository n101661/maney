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
		user.Put("/accounts/{oid:uint64}", s.UpdateAccount)
		user.Delete("/accounts/{oid:uint64}", s.DeleteAccount)
	}
	{ // user's categories
		user.Post("/categories")
		user.Get("/categories")
		user.Put("/categories/{oid:uint64}")
		user.Delete("/categories/{oid:uint64}")
	}
	{ // user's shops
		user.Post("/shops")
		user.Get("/shops")
		user.Put("/shops/{oid:uint64}")
		user.Delete("/shops/{oid:uint64}")
	}
	{ // user's fees
		user.Post("/fees")
		user.Get("/fees")
		user.Put("/fees/{oid:uint64}")
		user.Delete("/fees/{oid:uint64}")
	}
	{ // user's daily items
		user.Post("/daily-items")
		user.Get("/daily-items")
		user.Put("/daily-items/{oid:uint64}")
		user.Delete("/daily-items/{oid:uint64}")
		user.Get("/daily-items/export")
		user.Post("/daily-items/import")
	}
	{ // user's repeating items
		user.Post("/repeating-items")
		user.Get("/repeating-items")
		user.Put("/repeating-items/{oid:uint64}")
		user.Delete("/repeating-items/{oid:uint64}")
	}
}

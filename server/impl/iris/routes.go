package iris

func (s *Server) registerRoutes() {

	s.app.Post("/log-in", s.LogIn)
	s.app.Post("/log-out", s.LogOut)
	s.app.Post("/sign-up", s.SignUp)

	user := s.app.Party("/{user-id:string}")

	{ // user's config
		user.Put("/config")
		user.Get("/config")
	}
	{ // user's accounts
		user.Post("/accounts")
		user.Get("/accounts")
		user.Put("/accounts/{account-id:int64}")
		user.Delete("/accounts/{account-id:int64}")
	}
	{ // user's categories
		user.Post("/categories")
		user.Get("/categories")
		user.Put("/categories/{category-id:int64}")
		user.Delete("/categories/{category-id:int64}")
	}
	{ // user's shops
		user.Post("/shops")
		user.Get("/shops")
		user.Put("/shops/{shop-id:int64}")
		user.Delete("/shops/{shop-id:int64}")
	}
	{ // user's fees
		user.Post("/fees")
		user.Get("/fees")
		user.Put("/fees/{fee-id:int64}")
		user.Delete("/fees/{fee-id:int64}")
	}
	{ // user's daily items
		user.Post("/daily-items")
		user.Get("/daily-items")
		user.Put("/daily-items/{daily-item-id:int64}")
		user.Delete("/daily-items/{daily-item-id:int64}")
		user.Get("/daily-items/export")
		user.Post("/daily-items/import")
	}
	{ // user's repeating items
		user.Post("/repeating-items")
		user.Get("/repeating-items")
		user.Put("/repeating-items/{repeating-item-id:int64}")
		user.Delete("/repeating-items/{repeating-item-id:int64}")
	}
}

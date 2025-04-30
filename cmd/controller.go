package main

import (
	"github.com/n101661/maney/server/accounts"
	"github.com/n101661/maney/server/categories"
	"github.com/n101661/maney/server/impl/iris"
	"github.com/n101661/maney/server/shops"
	"github.com/n101661/maney/server/users"
)

func newIrisController(services *Services) *iris.Controllers {
	return &iris.Controllers{
		User:     users.NewIrisController(services.User),
		Account:  accounts.NewIrisController(services.Account),
		Category: categories.NewIrisController(services.Category),
		Shop:     shops.NewIrisController(services.Shop),
	}
}

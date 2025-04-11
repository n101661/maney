package middleware

import (
	"github.com/kataras/iris/v12"
)

func ErrorHandler(c iris.Context) {
	err := c.GetErr()
	if err != nil {
		c.Text("see request id [%v] for more details", c.GetID())
	}
}

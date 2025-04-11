package recover

import (
	"fmt"
	"runtime/debug"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func New() context.Handler {
	return func(ctx *context.Context) {
		defer func() {
			if recoveredMessage := recover(); recoveredMessage != nil {
				err := fmt.Errorf("%v\n%s", recoveredMessage, debug.Stack())
				ctx.StopWithError(iris.StatusInternalServerError, iris.PrivateError(err))
			}
		}()
		ctx.Next()
	}
}

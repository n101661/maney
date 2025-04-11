package logger

import (
	"fmt"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/samber/lo"

	"github.com/n101661/maney/pkg/utils"
)

func New(opts ...utils.Option[options]) context.Handler {
	o := utils.ApplyOptions(&options{}, opts)
	return func(c *context.Context) {
		start := time.Now()

		c.Next()

		statusCode := c.GetStatusCode()
		if logger := c.Application().Logger(); logger.Level == golog.DebugLevel || internalError(statusCode) {
			var message strings.Builder

			printRequest := true
			if o.excludeRequest != nil {
				printRequest = !o.excludeRequest(c)
			}

			requestID := ""
			if o.requestIDFunc != nil {
				requestID = o.requestIDFunc(c)
			}

			method := c.Method()

			message.WriteString(fmt.Sprintf(
				"[%s] %s %s returns %d status code in %d ms",
				requestID,
				method,
				lo.IfF(printRequest, func() string {
					return c.Request().URL.RequestURI()
				}).ElseF(func() string {
					return c.Path()
				}),
				c.GetStatusCode(),
				time.Since(start).Milliseconds(),
			))
			if method == "POST" || method == "PUT" {
				if printRequest {
					message.WriteByte('\n')

					raw, err := httputil.DumpRequest(c.Request(), true)
					if err != nil {
						raw = fmt.Appendf(raw, "<failed to extract the request: %v>", err)
					}
					message.WriteString(fmt.Sprintf("Request: %s", raw))
				}
			}
			if o.excludeResponseError == nil || !o.excludeResponseError(c) {
				err := c.GetErr()
				if err != nil {
					message.WriteByte('\n')
					message.WriteString(fmt.Sprintf("Error: %s", err.Error()))
				}
			}

			if internalError(statusCode) {
				logger.Warn(message.String())
			} else {
				logger.Debug(message.String())
			}
		}
	}
}

func internalError(code int) bool {
	return code >= 500 && code < 600
}

type options struct {
	requestIDFunc        func(ctx iris.Context) string
	excludeRequest       func(ctx iris.Context) bool
	excludeResponseError func(ctx iris.Context) bool
}

func WithRequestIDFunc(f func(ctx iris.Context) string) utils.Option[options] {
	return func(o *options) {
		o.requestIDFunc = f
	}
}

func WithExcludeRequest(f func(ctx iris.Context) bool) utils.Option[options] {
	return func(o *options) {
		o.excludeRequest = f
	}
}

func WithExcludeResponseError(f func(ctx iris.Context) bool) utils.Option[options] {
	return func(o *options) {
		o.excludeResponseError = f
	}
}

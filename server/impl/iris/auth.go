package iris

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/x/errors"
	"github.com/n101661/maney/pkg/services/auth"
)

func (s *Server) ValidateAccessToken(c *context.Context) {
	h := c.GetHeader(headerAuthorization)
	if h == "" {
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	frags := strings.SplitN(h, " ", 2)
	if frags[0] != authType || len(frags) != 2 {
		c.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	_, err := s.authService.ValidateAccessToken(c.Request().Context(), frags[1])
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrTokenExpired) {
			c.StopWithStatus(iris.StatusUnauthorized)
		} else {
			c.StopWithError(iris.StatusInternalServerError, err)
		}
		return
	}
}

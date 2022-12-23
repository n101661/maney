package iris

import (
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"

	"github.com/n101661/maney/server/models"
)

func (s *Server) LogIn(ctx iris.Context) {
	var r models.LogInBody
	if err := ctx.ReadJSON(&r); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	user, err := s.db.User().Get(r.Id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	if user == nil ||
		bcrypt.CompareHashAndPassword(user.Password, encryptPassword(r.Password)) != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	tokenMaxAge := 8 * time.Hour

	token, err := s.auth.generateToken(map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"expiry": time.Now().Add(tokenMaxAge),
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.Header(
		"Set-Cookie",
		fmt.Sprintf("token=%s; Max-Age=%d; HttpOnly", token, tokenMaxAge/time.Second),
	)
	ctx.StatusCode(iris.StatusOK)
}

func encryptPassword(password string) []byte {
	h := sha512.New()
	h.Write([]byte(password))
	return h.Sum(nil)
}

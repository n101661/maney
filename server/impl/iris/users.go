package iris

import (
	"errors"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"

	"github.com/n101661/maney/database"
	dbModels "github.com/n101661/maney/database/models"
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
		s.auth.validatePassword(user.Password, []byte(r.Password)) != nil {
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

func (s *Server) LogOut(ctx iris.Context) {
	ctx.Header("Set-Cookie", `token=""; Max-Age=0; HttpOnly`)
	ctx.StatusCode(iris.StatusOK)
}

func (s *Server) SignUp(ctx iris.Context) {
	var r models.SignUpBody
	if err := ctx.ReadJSON(&r); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	pwd, err := s.auth.encryptPassword(r.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	err = s.db.User().Create(dbModels.User{
		ID:       r.Id,
		Name:     r.Name,
		Password: pwd,
	})
	if err != nil {
		if errors.Is(err, database.ErrResourceExisted) {
			ctx.StatusCode(iris.StatusConflict)
			ctx.WriteString("the user id has existed")
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
}

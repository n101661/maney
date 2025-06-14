package iris

import (
	"context"

	"github.com/kataras/iris/v12"

	"github.com/n101661/maney/server/models"
)

type SimpleCreateTemplate[RequestBody, ServiceRequest, ServiceReply, ResponseBody any] struct {
	Service interface {
		Create(context.Context, *ServiceRequest) (*ServiceReply, error)
	}

	// ParseServiceRequest the returned error is considered as user bad request and write 400 status code.
	// If you want to write 500 status code, wrap the error by InternalError function.
	ParseServiceRequest func(userID string, r *RequestBody) (*ServiceRequest, error)
	// BadRequest checks if the error returned from Service is http bad request or not.
	BadRequest       func(err error) (httpCode int, yes bool)
	ParseAPIResponse func(*ServiceReply) (*ResponseBody, error)
}

func (t *SimpleCreateTemplate[RequestBody, ServiceRequest, ServiceReply, ResponseBody]) Create(c iris.Context) {
	var r RequestBody
	if err := c.ReadJSON(&r); err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	sr, err := t.ParseServiceRequest(userID, &r)
	if err != nil {
		if e, ok := err.(*internalError); ok {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(e.err))
		} else {
			c.StopWithText(iris.StatusBadRequest, err.Error())
		}
		return
	}

	reply, err := t.Service.Create(c.Request().Context(), sr)
	if err != nil {
		if code, y := t.BadRequest(err); y {
			c.StopWithText(code, err.Error())
			return
		}
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	resp, err := t.ParseAPIResponse(reply)
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, resp)
}

type SimpleListTemplate[ServiceRequest, ServiceReply, ResponseBody any] struct {
	Service interface {
		List(context.Context, *ServiceRequest) (*ServiceReply, error)
	}

	// ParseServiceRequest the returned error is considered as user bad request and write 400 status code.
	// If you want to write 500 status code, wrap the error by InternalError function.
	ParseServiceRequest func(c iris.Context, userID string) (*ServiceRequest, error)
	// BadRequest checks if the error returned from Service is http bad request or not.
	BadRequest       func(err error) (httpCode int, yes bool)
	ParseAPIResponse func(*ServiceReply) (*ResponseBody, error)
}

func (t *SimpleListTemplate[ServiceRequest, ServiceReply, ResponseBody]) List(c iris.Context) {
	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	sr, err := t.ParseServiceRequest(c, userID)
	if err != nil {
		if e, ok := err.(*internalError); ok {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(e.err))
		} else {
			c.StopWithText(iris.StatusBadRequest, err.Error())
		}
	}

	reply, err := t.Service.List(c.Request().Context(), sr)
	if err != nil {
		if code, y := t.BadRequest(err); y {
			c.StopWithText(code, err.Error())
			return
		}
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	resp, err := t.ParseAPIResponse(reply)
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, resp)
}

type SimpleUpdateTemplate[RequestBody, ServiceRequest, ServiceReply any] struct {
	// Placeholder is the ID of the placeholder in API path.
	Placeholder string
	Service     interface {
		Update(context.Context, *ServiceRequest) (*ServiceReply, error)
	}

	// ParseServiceRequest the returned error is considered as user bad request and write 400 status code.
	// If you want to write 500 status code, wrap the error by InternalError function.
	ParseServiceRequest func(userID string, publicID string, r *RequestBody) (*ServiceRequest, error)
	// BadRequest checks if the error returned from Service is http bad request or not.
	BadRequest func(err error) (httpCode int, yes bool)
}

func (t *SimpleUpdateTemplate[RequestBody, ServiceRequest, ServiceReply]) Update(c iris.Context) {
	publicID := c.Params().GetString(t.Placeholder)

	var r RequestBody
	if err := c.ReadJSON(&r); err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	sr, err := t.ParseServiceRequest(userID, publicID, &r)
	if err != nil {
		if e, ok := err.(*internalError); ok {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(e.err))
		} else {
			c.StopWithText(iris.StatusBadRequest, err.Error())
		}
		return
	}

	_, err = t.Service.Update(c.Request().Context(), sr)
	if err != nil {
		if code, y := t.BadRequest(err); y {
			c.StopWithText(code, err.Error())
			return
		}
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.EmptyResponse{})
}

type SimpleDeleteTemplate[ServiceRequest, ServiceReply any] struct {
	// Placeholder is the ID of the placeholder in API path.
	Placeholder string
	Service     interface {
		Delete(context.Context, *ServiceRequest) (*ServiceReply, error)
	}

	ParseServiceRequest func(userID string, publicID string) *ServiceRequest
	// BadRequest checks if the error returned from Service is http bad request or not.
	BadRequest func(err error) (httpCode int, yes bool)
}

func (t *SimpleDeleteTemplate[ServiceRequest, ServiceReply]) Delete(c iris.Context) {
	publicID := c.Params().GetString(t.Placeholder)

	user := c.User()
	if user == nil {
		c.StopWithJSON(iris.StatusUnauthorized, &models.EmptyResponse{})
		return
	}

	userID, err := user.GetID()
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	sr := t.ParseServiceRequest(userID, publicID)

	_, err = t.Service.Delete(c.Request().Context(), sr)
	if err != nil {
		if code, y := t.BadRequest(err); y {
			c.StopWithText(code, err.Error())
			return
		}
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.EmptyResponse{})
}

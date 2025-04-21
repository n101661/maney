package categories

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/samber/lo"

	"github.com/n101661/maney/server/models"
	"github.com/n101661/maney/server/repository"
)

type IrisController struct {
	s Service
}

func NewIrisController(s Service) *IrisController {
	return &IrisController{
		s: s,
	}
}

func (controller *IrisController) Create(c iris.Context) {
	var r models.CreatingCategory
	if err := c.ReadJSON(&r); err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	type_, err := parseType(string(r.Type))
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, err.Error())
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

	reply, err := controller.s.Create(c.Request().Context(), &CreateRequest{
		UserID: userID,
		Type:   type_,
		Category: &BaseCategory{
			Name:   r.Name,
			IconID: int32(lo.FromPtrOr(r.IconId, 0)),
		},
	})
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.ObjectId{
		Id: lo.ToPtr(models.Id(reply.Category.ID)),
	})
}

func (controller *IrisController) List(c iris.Context) {
	type_, err := parseType(c.URLParamDefault("type", repository.CategoryTypeExpense.String()))
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, err.Error())
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

	reply, err := controller.s.List(c.Request().Context(), &ListRequest{
		UserID: userID,
		Type:   type_,
	})
	if err != nil {
		c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		return
	}

	c.StopWithJSON(iris.StatusOK, lo.Map(reply.Categories, func(item *Category, _ int) *models.Category {
		return &models.Category{
			Id:     lo.ToPtr(models.Id(item.ID)),
			Name:   item.Name,
			IconId: lo.ToPtr(models.Id(item.IconID)),
		}
	}))
}

func (controller *IrisController) Update(c iris.Context) {
	categoryID, err := c.Params().GetInt32("categoryId")
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, "missing or invalid category id")
		return
	}

	var r models.BasicCategory
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

	_, err = controller.s.Update(c.Request().Context(), &UpdateRequest{
		UserID:     userID,
		CategoryID: categoryID,
		Category: &BaseCategory{
			Name:   r.Name,
			IconID: int32(lo.FromPtrOr(r.IconId, 0)),
		},
	})
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.StopWithText(iris.StatusNotFound, "no [%d] category id", categoryID)
		} else {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		}
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.EmptyResponse{})
}

func (controller *IrisController) Delete(c iris.Context) {
	categoryID, err := c.Params().GetInt32("categoryId")
	if err != nil {
		c.StopWithText(iris.StatusBadRequest, "missing or invalid category id")
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

	_, err = controller.s.Delete(c.Request().Context(), &DeleteRequest{
		UserID:     userID,
		CategoryID: categoryID,
	})
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.StopWithText(iris.StatusNotFound, "no [%d] category id", categoryID)
		} else {
			c.StopWithPlainError(iris.StatusInternalServerError, iris.PrivateError(err))
		}
		return
	}

	c.StopWithJSON(iris.StatusOK, &models.EmptyResponse{})
}

func parseType(s string) (Type, error) {
	return repository.ToCategoryType(s)
}

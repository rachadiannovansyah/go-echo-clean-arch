package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/rachadiannovansyah/go-echo-clean-arch/domain"
	errHandle "github.com/rachadiannovansyah/go-echo-clean-arch/utils"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// UserHandler  represent the httphandler for article
type UserHandler struct {
	UserUcase domain.UserUsecase
}

// NewUserHandler will initialize the articles/ resources endpoint
func NewUserHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		UserUcase: us,
	}
	e.GET("/users", handler.FetchUser)
}

// FetchUser will fetch the article based on given params
func (a *UserHandler) FetchUser(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listUser, nextCursor, err := a.UserUcase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	fmt.Println(nextCursor)
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listUser)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case errHandle.ErrInternalServerError:
		return http.StatusInternalServerError
	case errHandle.ErrNotFound:
		return http.StatusNotFound
	case errHandle.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

package routers

import (
	"main/core"
	"main/core/business"
	"main/core/models"
	"net/http"
	"strconv"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

type UserRouter struct {
	Name string
	g    *echo.Group
}

func (r *UserRouter) Connect(s *core.Server) {
	r.g = s.Echo.Group(r.Name)

	user := business.UserBusiness{
		DB: s.DB,
	}

	r.g.GET("/", func(c echo.Context) error {
		token, _ := c.Get("user").(auth.Token)
		return c.JSON(http.StatusOK, token.Firebase.Identities)
	}, s.AuthMiddleware)

	r.g.POST("/", func(c echo.Context) (err error) {
		profile := new(models.UserProfile)

		if err = c.Bind(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = user.Create(*profile)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, profile)
	})

	// [PUT] Input: Body (UserProfileUpdated)
	r.g.PUT("/:id", func(c echo.Context) (err error) {
		profile := new(models.UserProfileUpdated)
		id, _ := strconv.Atoi(c.Param("id"))
		if err = c.Bind(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = user.Update(strconv.Itoa((id)), *profile)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, id)
	})
	// [DELETE] Input user's id
	r.g.DELETE("/:id", func(c echo.Context) (err error) {
		id, _ := strconv.Atoi(c.Param("id"))

		err = user.Delete(strconv.Itoa(id))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, id)
	})
}

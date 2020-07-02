package main

import (
	"github.com/s-tajima/nspv"

	"github.com/labstack/echo/v4"
	"net/http"
)

type User struct {
	Passowrd string `form:"password"`
}

func main() {
	e := echo.New()

	e.POST("/password", postPassword)

	e.Logger.Fatal(e.Start(":1323"))
}

func postPassword(c echo.Context) error {
	v := nspv.NewValidator()

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	res, _ := v.Validate(u.Passowrd)
	return c.String(http.StatusOK, res.String())
}

package middlewares

import (
	"go-admin/util"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	if _, err := util.ParseJwt(cookie); err != nil {
		c.Status(http.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	return c.Next()

}

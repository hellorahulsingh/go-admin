package controllers

import (
	"go-admin/database"
	"go-admin/models"
	"go-admin/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["confirm_password"] {
		c.Status(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Passwords did not match",
		})
	}

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    1,
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := user.ComparePasswords(data["password"]); err != nil {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	tokenstring, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))

	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenstring,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"token": tokenstring,
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	if user.Id == 0 {
		c.Status(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"token": "Logged out successfully",
	})
}

package controllers

import (
	"go-admin/database"
	"go-admin/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const Password string = "12345"

func PaginatedUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	// offset := (page - 1) * limit

	// var total int64
	// var users []models.User

	// database.DB.Preload("Role").Offset(offset).Limit(limit).Find(&users)
	// database.DB.Model(&models.User{}).Count(&total)

	// return c.JSON(fiber.Map{
	// 	"data": users,
	// 	"meta": fiber.Map{
	// 		"page":       page,
	// 		"total":      total,
	// 		"limit":      limit,
	// 		"total_page": float64(int(total) / limit),
	// 	},
	// })
	return c.JSON(models.Paginate(database.DB, &models.User{}, page))
}

func AllUsers(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Preload("Role").Find(&users)

	return c.JSON(users)
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	user.SetPassword(Password)

	database.DB.Create(&user)

	return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	// FIRST WAY
	// var user models.User
	// database.DB.Where("id = ?", uint(id)).First(&user)
	// if user.Id == 0 {
	// 	c.Status(http.StatusNotFound)
	// 	return c.JSON(fiber.Map{
	// 		"message": "User not found",
	// 	})
	// }

	// SECOND WAY
	var user = models.User{
		Id: uint(id),
	}
	database.DB.Preload("Role").Find(&user)

	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	database.DB.Delete(&user)

	return nil
}

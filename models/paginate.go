package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const limit int = 5

func Paginate(db *gorm.DB, entity Entity, page int) fiber.Map {
	offset := (page - 1) * limit

	data := entity.Take(db, limit, offset)
	count := entity.Count(db)

	return fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":       page,
			"total":      count,
			"limit":      limit,
			"total_page": float64(int(count) / limit),
		},
	}
}

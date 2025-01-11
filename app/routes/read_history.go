package routes

import (
	"bible-users-api/app/appdata"
	"bible-users-api/app/models"
	"bible-users-api/app/utils"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MarkChapterAsRead(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	book := c.FormValue("book")
	chapterStr := c.FormValue("chapter")
	chapterInt, err := strconv.Atoi(chapterStr)
	chapter := uint(chapterInt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chapter value",
		})
	}
	var bookStruct appdata.Book
	for _, b := range appdata.Books {
		if b.Book == book {
			bookStruct = b
			break
		}
	}
	if bookStruct.Book != book {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bible book",
		})
	}
	if chapter > bookStruct.Chapters {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chapter number",
		})
	}
	readHistory := models.ReadHistory{UserID: user_id, Book: book, Chapter: chapter}
	result := appdata.DB.Create(&readHistory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			// Successful response is intentional
			return c.JSON(fiber.Map{
				"error": "This chapter is already marked read",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong, try again later",
			})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Chapter marked read",
	})
}

func MarkChapterAsUnread(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	book := c.FormValue("book")
	chapterStr := c.FormValue("chapter")
	chapterInt, err := strconv.Atoi(chapterStr)
	chapter := uint(chapterInt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chapter value",
		})
	}
	var readHistory models.ReadHistory
	result := appdata.DB.Where("user_id = ? AND book = ? AND chapter = ?", user_id, book, chapter).First(&readHistory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(fiber.Map{
				"error": "This book and chapter is not marked as read",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong, try again later",
			})
		}
	}
	appdata.DB.Delete(&readHistory)
	return c.JSON(fiber.Map{
		"message": "Chapter marked unread",
	})
}

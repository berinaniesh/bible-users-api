package routes

import (
	"bible-users-api/app/appdata"
	"bible-users-api/app/models"
	"bible-users-api/app/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddBookmark(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	book := c.FormValue("book")
	var bookStruct appdata.Book
	for _, b := range appdata.Books {
		if b.Book == book {
			bookStruct = b
			break
		}
	}
	if bookStruct.Book == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Book not valid",
		})
	}
	chapterString := c.FormValue("chapter")
	verseNumberString := c.FormValue("verse")
	chapterInt, err := strconv.Atoi(chapterString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chapter number not a valid number",
		})
	}
	verseNumberInt, err := strconv.Atoi(verseNumberString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Verse Number not a valid number",
		})
	}
	var bookmark models.Bookmark = models.Bookmark{UserID: user_id, Book: book, ChapterNumber: uint(chapterInt), VerseNumber: uint(verseNumberInt)}
	appdata.DB.Create(&bookmark)
	return c.JSON(fiber.Map{
		"message": "Created bookmark",
	})
}

func DeleteBookmark(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	book := c.FormValue("book")
	chapterString := c.FormValue("chapter")
	verseNumberString := c.FormValue("verse")
	chapterInt, err := strconv.Atoi(chapterString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chapter number not a valid number",
		})
	}
	verseNumberInt, err := strconv.Atoi(verseNumberString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Verse Number not a valid number",
		})
	}
	var bookmark models.Bookmark
	appdata.DB.Where("user_id = ? AND book = ? AND chapter_number = ? and verse_number = ?", user_id, book, uint(chapterInt), uint(verseNumberInt)).First(&bookmark)
	appdata.DB.Delete(&bookmark)
	return c.JSON(fiber.Map{
		"message": "Bookmark deleted",
	})
}

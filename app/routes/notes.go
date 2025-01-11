package routes

import (
	"bible-users-api/app/appdata"
	"bible-users-api/app/models"
	"bible-users-api/app/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateNote(c *fiber.Ctx) error {
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
	noteString := c.FormValue("note")
	if noteString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Note empty",
		})
	}
	note := models.Note{UserID: user_id, Book: book, ChapterNumber: uint(chapterInt), VerseNumber: uint(verseNumberInt), Note: noteString}
	appdata.DB.Create(&note)
	return c.JSON(note)
}

func DeleteNote(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	noteIdString := c.Params("noteid")
	noteId, err := strconv.Atoi(noteIdString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wrong note id",
		})
	}
	var note models.Note
	appdata.DB.First(&note, noteId)
	if note.UserID != user_id {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Note not found or it doesn't belong to you.",
		})
	}
	appdata.DB.Delete(&note)
	return c.JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}

func UpdateNote(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	noteIdString := c.Params("noteid")
	noteId, err := strconv.Atoi(noteIdString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Wrong note id",
		})
	}
	var note models.Note
	appdata.DB.First(&note, noteId)
	if note.UserID != user_id {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Note not found or it doesn't belong to you.",
		})
	}
	noteString := c.FormValue("note")
	if noteString != "" {
		note.Note = noteString
		appdata.DB.Save(&note)
	}
	return c.JSON(note)
}

func GetNotesOfUser(c *fiber.Ctx) error {
	user_id := utils.GetUserFromJwt(c)
	book := c.Query("book") // prefer using abbreviation (there are spaces in book names)
	abbreviation := c.Query("abbreviation")
	var notes []models.Note
	chapterString := c.Query("chapter")
	if book == "" && abbreviation == "" {
		appdata.DB.Where("user_id = ?", user_id).Find(&notes)
		return c.JSON(notes)
	}
	chapterInt, _ := strconv.Atoi(chapterString)
	if book != "" {
		if chapterInt != 0 {
			appdata.DB.Where("user_id = ? AND book = ? AND chapter_number = ?", user_id, book, uint(chapterInt)).Find(&notes)
			return c.JSON(notes)
		} else {
			appdata.DB.Where("user_id = ? AND book = ?", user_id, book).Find(&notes)
			return c.JSON(notes)
		}
	}
	if abbreviation != "" {
		for _, b := range appdata.Books {
			if b.Abbreviation == abbreviation {
				if chapterInt != 0 {
					appdata.DB.Where("user_id = ? AND book = ? AND chapter_number = ?", user_id, b.Book, uint(chapterInt)).Find(&notes)
					return c.JSON(notes)
				} else {
					appdata.DB.Where("user_id = ? AND book = ?", user_id, b.Book).Find(&notes)
					return c.JSON(notes)
				}
			}
		}
	}
	return c.JSON(notes)
}

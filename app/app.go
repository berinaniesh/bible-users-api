package app

import (
	"bible-users-api/app/appdata"
	"bible-users-api/app/models"
	"bible-users-api/app/routes"
	"log"
	"os"
	"strconv"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Fiber *fiber.App
}

func NewApp() *App {
	fiberApp := fiber.New(fiber.Config{
		AppName: "Bible Users API",
	})
	return &App{
		Fiber: fiberApp,
	}
}

func getExpiryMinutes(s string) uint {
	value := os.Getenv(s)
	if value == "" {
		log.Fatal(s + " is not set")
	}
	valueUint, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		log.Fatal(s + " could not be parsed as minutes")
	}
	return uint(valueUint)
}

func (app *App) InitializeApp() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment variables from .env file")
	}
	appdata.JwtExpiryMinutes = getExpiryMinutes("JWT_EXPIRY_MINUTES")
	appdata.RefreshExpiryMinutes = getExpiryMinutes("REFRESH_EXPIRY_MINUTES")
	appdata.RefreshExpiryNoRemember = getExpiryMinutes("REFRESH_EXPIRY_NO_REMEMBER")
	appdata.JwtExpiryNoRemember = getExpiryMinutes("JWT_EXPIRY_NO_REMEMBER")
	appdata.SmtpServer = os.Getenv("SMTP_SERVER")
	appdata.SmtpPassword = os.Getenv("SMTP_PASSWORD")
	envSmtpPort, _ := strconv.ParseUint(os.Getenv("SMTP_PORT"), 10, 32)
	appdata.SmtpPort = uint(envSmtpPort)
	appdata.SmtpUsername = os.Getenv("SMTP_FROM")
	if appdata.SmtpPort == 0 || appdata.SmtpServer == "" || appdata.SmtpPassword == "" || appdata.SmtpUsername == "" {
		log.Fatal("Failed to load environment variables for SMTP settings.")
	}
	jwtSecretString := os.Getenv("JWT_SECRET")
	if jwtSecretString == "" {
		log.Fatal("Failed to load JWT_SECRET from .env")
	}
	appdata.JwtSecret = []byte(jwtSecretString)
	appdata.ResetValidMinutes = getExpiryMinutes("RESET_VALID_MINUTES")
}

func (app *App) InitializeDatabase() {
	var err error
	dsn := os.Getenv("DSN")
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)
	appdata.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         gormLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.RefreshToken{},
		&models.ForgotPassword{},
		&models.VerifyEmail{},
		&models.ReadHistory{},
		&models.UserPreference{},
		&models.Bookmark{},
		&models.Note{},
	}
	for _, model := range modelsToMigrate {
		if err := appdata.DB.AutoMigrate(model); err != nil {
			log.Fatal("Failed to do database migrations")
		}
	}
}

func (app *App) SetupRoutes() {
	app.Fiber.Use(recover.New())
	app.Fiber.Get("/", routes.Home)
	app.Fiber.Post("/users", routes.CreateUser)
	app.Fiber.Post("/login", routes.LoginUser)
	app.Fiber.Post("/refreshtoken", routes.RefreshToken)
	app.Fiber.Post("/sendforgotpasswordemail", routes.SendForgotPasswordEmail)
	app.Fiber.Post("/resetpassword", routes.ResetPassword)
	app.Fiber.Post("/verifyemail", routes.VerifyEmail)
	app.Fiber.Post("/logout", routes.Logout)

	app.Fiber.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: appdata.JwtSecret},
	}))

	app.Fiber.Post("/sendemailverificationemail", routes.SendEmailVerificationEmail)
	app.Fiber.Put("/users", routes.UpdateUser)
	app.Fiber.Get("/me", routes.GetSelfInfo)
	app.Fiber.Post("/logoutall", routes.LogoutAll)
	app.Fiber.Post("/changepassword", routes.ChangePassword)
	app.Fiber.Post("/markchapterasread", routes.MarkChapterAsRead)
	app.Fiber.Delete("/markchapterasread", routes.MarkChapterAsUnread)
	app.Fiber.Put("/userpreferences", routes.UpdateUserPreferences)
	app.Fiber.Delete("/userpreferences", routes.DeleteUserPreferences)
	app.Fiber.Post("/bookmark", routes.AddBookmark)
	app.Fiber.Delete("/bookmark", routes.DeleteBookmark)
	app.Fiber.Post("/note", routes.CreateNote)
	app.Fiber.Delete("/note/:noteid", routes.DeleteNote)
	app.Fiber.Put("/note/:noteid", routes.UpdateNote)
	app.Fiber.Get("/note", routes.GetNotesOfUser)
}

func (app *App) Start() {
	hostUrl := os.Getenv("HOST_URL")
	log.Fatal(app.Fiber.Listen(hostUrl))
}

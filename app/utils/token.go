package utils

import (
	"bible-users-api/app/appdata"
	"bible-users-api/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func PrepareAccessToken(user *models.User, remember bool) string {
	var jwtExpiryMinutes uint
	if remember {
		jwtExpiryMinutes = appdata.JwtExpiryMinutes
	} else {
		jwtExpiryMinutes = appdata.JwtExpiryNoRemember
	}
	expiry := time.Now().Add(time.Duration(jwtExpiryMinutes) * time.Minute)
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": expiry.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(appdata.JwtSecret)
	if err == nil {
		return signedToken
	} else {
		return ""
	}
}

func PrepareRefreshToken(user *models.User, device *string, location *string, remember bool) string {
	var refreshExpiryMinutes uint
	if remember {
		refreshExpiryMinutes = appdata.RefreshExpiryMinutes
	} else {
		refreshExpiryMinutes = appdata.RefreshExpiryNoRemember
	}
	expiry := time.Now().Add(time.Duration(refreshExpiryMinutes) * time.Minute)

	oneRefreshPeriodBefore := time.Now().Add(-time.Duration(appdata.RefreshExpiryMinutes) * time.Minute)
	appdata.DB.Where("expires_at < ?", oneRefreshPeriodBefore).Delete(&models.RefreshToken{})
	tokenString := GenerateAlphanumeric(25)
	refreshToken := models.RefreshToken{UserID: user.ID, Device: *device, Location: location, Token: tokenString, ExpiresAt: expiry, Remember: remember}
	result := appdata.DB.Create(&refreshToken)
	if result.Error == nil {
		return tokenString
	} else {
		return ""
	}
}

func GetUserFromJwt(c *fiber.Ctx) uint {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	user_id := uint(claims["id"].(float64))
	return user_id
}

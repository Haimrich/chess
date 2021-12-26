package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

var key = os.Getenv("JWT_SECRET_KEY")

func GenerateTokensAndSetCookies(uid string, c *gin.Context) (string, string, error) {

	accessTokenExpiration := time.Now().Add(15 * time.Minute)
	accessToken, err := generateToken(uid, accessTokenExpiration, key)
	if err != nil {
		return "", "", err
	}

	refreshTokenExpiration := time.Now().Add(24 * time.Hour)
	refreshToken, err := generateToken(uid, refreshTokenExpiration, key)
	if err != nil {
		return "", "", err
	}

	setCookie(accessTokenCookieName, accessToken, accessTokenExpiration, c)
	setCookie(refreshTokenCookieName, refreshToken, refreshTokenExpiration, c)

	return accessToken, refreshToken, nil
}

func DeleteTokens(c *gin.Context) {
	deleteCookie(accessTokenCookieName, c)
	deleteCookie(refreshTokenCookieName, c)
}

func generateToken(uid string, expirationTime time.Time, secret string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTime.Unix()
	claims["uid"] = uid

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func setCookie(name string, value string, expiration time.Time, c *gin.Context) {
	c.SetCookie(name, value, expiration.Second(), "/", "", false, true)
}

func deleteCookie(name string, c *gin.Context) {
	setCookie(name, "", time.Unix(0, 0), c)
}

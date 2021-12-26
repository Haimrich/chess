package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(authNeeded bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		accessTokenString, refreshTokenString, err := extractTokens(c)
		if err != nil {
			fmt.Println("[AUTH] Cookies not found.")
			if authNeeded {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Abort()
			} else {
				c.Next()
			}
			return
		}

		accessToken, err := jwt.Parse(accessTokenString, keyFunc)

		if err == nil && accessToken.Valid {

			accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
			fmt.Printf("[AUTH] Access token valido per %s.\n", accessTokenClaims["uid"])
			c.Set("uid", accessTokenClaims["uid"])
			c.Next()

		} else {

			fmt.Printf("[AUTH] Access %s. Trying refresh token.\n", err)
			refreshToken, err := jwt.Parse(refreshTokenString, keyFunc)

			if err != nil || !refreshToken.Valid {
				fmt.Printf("[AUTH] Refresh %s :(\n", err)
				if authNeeded {
					c.JSON(http.StatusUnauthorized, err.Error())
					c.Abort()
				} else {
					c.Next()
				}

			} else {

				uid := refreshToken.Claims.(jwt.MapClaims)["uid"].(string)
				fmt.Printf("[AUTH] Trovato refresh token valido per %s.\n", uid)
				refreshTokens(uid, c)
				c.Set("uid", uid)
				c.Next()

			}
		}

	}
}

func extractTokens(c *gin.Context) (string, string, error) {
	accessCookie, err := c.Request.Cookie(accessTokenCookieName)
	if err != nil {
		return "", "", err
	}
	refreshCookie, err := c.Request.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", "", err
	}

	return accessCookie.Value, refreshCookie.Value, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(key), nil
}

func refreshTokens(uid string, c *gin.Context) error {
	fmt.Printf("[AUTH] Refreshing tokens for %s.\n", uid)
	_, _, err := GenerateTokensAndSetCookies(uid, c)
	return err
}

/*
	if time.Until(time.Unix(accessTokenClaims.ExpiresAt, 0)) < 15*time.Minute {
		refreshTokens(accessTokenClaims.uid, c)
	}
*/

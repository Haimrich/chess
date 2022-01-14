package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const logPrefix = "      Auth - "

func TokenAuthMiddleware(authNeeded bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		accessTokenCookie, err := c.Request.Cookie(accessTokenCookieName)
		if err == nil {
			accessTokenString := accessTokenCookie.Value

			var accessToken *jwt.Token
			accessToken, err = jwt.Parse(accessTokenString, keyFunc)

			if err == nil && accessToken.Valid {
				accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
				fmt.Printf(logPrefix+"Access token valido per %s.\n", accessTokenClaims["uid"])
				c.Set("uid", accessTokenClaims["uid"])
				c.Next()
				return
			}
		}

		fmt.Printf(logPrefix+"Access %s. Trying refresh token.\n", err)

		refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
		if err == nil {
			refreshTokenString := refreshTokenCookie.Value

			var refreshToken *jwt.Token
			refreshToken, err = jwt.Parse(refreshTokenString, keyFunc)

			if err == nil && refreshToken.Valid {
				uid := refreshToken.Claims.(jwt.MapClaims)["uid"].(string)
				fmt.Printf(logPrefix+"Trovato refresh token valido per %s.\n", uid)
				refreshTokens(uid, c)
				c.Set("uid", uid)
				c.Next()
				return
			}
		}

		fmt.Printf(logPrefix+"Refresh %s :(\n", err)

		if authNeeded {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(key), nil
}

func refreshTokens(uid string, c *gin.Context) error {
	fmt.Printf(logPrefix+"Refreshing tokens for %s.\n", uid)
	_, _, err := GenerateTokensAndSetCookies(uid, c)
	return err
}

/*
	if time.Until(time.Unix(accessTokenClaims.ExpiresAt, 0)) < 15*time.Minute {
		refreshTokens(accessTokenClaims.uid, c)
	}
*/

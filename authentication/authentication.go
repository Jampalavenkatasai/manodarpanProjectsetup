package Authentication

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	//"manodarpanNewproject/server"
	"net/http"
	"strings"
	"time"
)

type JWTClaim struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("supersecretkeyvdjwbdhwjdbiwuhdqwihdiq")

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Writer.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "*")
		header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		header.Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.Writer.WriteHeader(http.StatusOK)
			return
		}

		c.Next()
		return
	}
}

func GenerateJWT(userId uint) (tokenString string, err error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims := &JWTClaim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "manodarpan",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		tokenString, tokenerror := StripBearerPrefixFromTokenString(tokenString)

		if tokenerror != nil {
			context.JSON(http.StatusGone, gin.H{"error": "error while parsing the authorization token"})
			context.Abort()
			return
		}
		fmt.Println("got jwt ", tokenString)
		//from
		userId, err := ValidateToken(tokenString)
		if err != nil {
			context.JSON(http.StatusGone, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		//result, err := s.DAO.CheckUserID(c, fmt.Sprint(userId))
		//if err != nil {
		//	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		//	return
		//}

		context.Set("currentUser", userId)

		context.Next()
	}
}
func StripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

// ValidateToken validates the JWT token and returns user data if successful.
func ValidateToken(signedToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		log.Fatal("error occurred during parsing the token", err.Error())
		return 0, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		log.Fatal("error occurred during parsing the token", err.Error())
		return 0, err
	}

	if token.Valid {
		return claims.UserId, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		err = errors.New("that's not even a token")
		log.Warning("That's not even a token", err.Error())
		return 0, err
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		err = errors.New("token is either expired or not active yet")
		log.Warning("token is either expired or not active yet", err.Error())
		return 0, err
	} else {
		err = errors.New("couldn't handle this token")
		log.Warning("Couldn't handle this token:", err.Error())
		return 0, err
	}
}

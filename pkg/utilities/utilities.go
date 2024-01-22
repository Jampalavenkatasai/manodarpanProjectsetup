/*
*******

package for common utilities

*********
*/
package utilities

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

const SUCCESS = "success"
const FAILED = "failed"
const IsRequired = "is required"
const jwtSecret = "secret"

func GetDuration(startTime time.Time) int {
	return int(time.Now().Sub(startTime).Round(time.Millisecond) / time.Millisecond)
}

func GetTotalPages(dataLimit int64, dataCount int64) int64 {

	totalPages := int64(math.Ceil(float64(dataCount) / float64(dataLimit)))
	return totalPages
}

func GetOffset(limit int64, page int) int {
	return (page - 1) * int(limit)
}

func GetLimit(ctx *gin.Context) int64 {
	limit1 := "20"
	if ctx.Request.PostFormValue("limit") != "" && ctx.Request.PostFormValue("limit") != "0" {
		limit1 = ctx.Request.PostFormValue("limit")
	}

	limit, _ := strconv.Atoi(limit1)

	return int64(limit)
}
func GetPage(ctx *gin.Context) (page int) {
	page1 := "1"
	if ctx.Request.PostFormValue("page") != "" && ctx.Request.PostFormValue("page") != "0" {
		page1 = ctx.Request.PostFormValue("page")
	}
	page, _ = strconv.Atoi(page1)
	return
}

func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return fmt.Sprintf("%s ", runtime.FuncForPC(pc).Name())
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Compare Password
func CheckPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
func SavePDF(file *multipart.FileHeader, uploadDir string, filename string) error {
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	return err
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Define the character set for the user ID
func GenerateUUID(prefix string, length int) string {
	currentTime := time.Now()
	year := currentTime.Year() % 100 // Get last two digits of the year
	month := int(currentTime.Month())
	rand.Seed(time.Now().UnixNano())
	var randomCode string
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		randomCode += string(charset[randomIndex])
	}
	newPrefix := fmt.Sprintf("Y%02dM%02d", year, month)
	return prefix + newPrefix + randomCode
}
func GetMissingFields(ctx *gin.Context) []string {
	requiredFields := []string{"first_name", "last_name", "email", "phone_no", "password"}
	missingFields := make([]string, 0)

	for _, field := range requiredFields {
		if ctx.PostForm(field) == "" {
			missingFields = append(missingFields, field)
		}
	}

	return missingFields
}

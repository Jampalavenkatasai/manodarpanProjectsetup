package server

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"manodarpanNewproject/authentication"
	"manodarpanNewproject/pkg/dbase"
	"manodarpanNewproject/pkg/models"
	"manodarpanNewproject/pkg/userdao"
	"manodarpanNewproject/pkg/utilities"
	"net/http"
	"strings"
	"time"
)

type Envs struct {
	Database  dbase.Config
	Port      string
	PageLimit int64
}
type Server struct {
	Env Envs
	DAO userdao.DAO
}

func (s *Server) UserRegistration(ctx *gin.Context) {
	// Parse form data
	firstName := ctx.PostForm("first_name")
	lastName := ctx.PostForm("last_name")
	email := ctx.PostForm("email")
	phoneNo := ctx.PostForm("phone_no")
	password := ctx.PostForm("password")

	// Validate mandatory fields
	missingFields := utilities.GetMissingFields(ctx)
	if len(missingFields) > 0 {
		errorMessage := "Mandatory fields are required: " + strings.Join(missingFields, ", ")
		ctx.JSON(http.StatusBadRequest, gin.H{"status": false, "message": errorMessage})
		return
	}

	// Check duplicate email
	existingUser, err := s.DAO.GetUserByEmail(email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("error", err)
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": false, "message": "Email is already in use"})
			return
		}
	}
	if existingUser.ID > 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"status": false, "message": "Email is already in use"})
		return
	}

	// Hash password

	hashedPassword, err := utilities.HashPassword(password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Error hashing password"})
		return
	}

	now := time.Now()
	//sql.NullString{String: regularString, Valid: true}
	// Create new user
	newAdminUser := models.User{
		FirstName: firstName,
		//LastName:  //sql.NullString{String: lastName, Valid: true},
		LastName: lastName,
		Email:    strings.ToLower(email),
		Password: hashedPassword,
		//PhoneNo:   sql.NullString{String: phoneNo, Valid: true},
		PhoneNo:   phoneNo,
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert user into the database
	result, err := s.DAO.CreateUser(ctx, newAdminUser)
	fmt.Println("Result1", result)
	if err != nil {
		log.Println("error1", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error while creating the User"})
		return
	}
	fmt.Println("Result2", result)

	if result.ID <= 0 {
		log.Println("error@result", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error while creating the User"})
		return
	}

	// Generate JWT token
	token, err := Authentication.GenerateJWT(uint(result.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error generating JWT token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "token": token, "message": "User registered successfully", "data": result})
}
func (s *Server) LoginUser(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	// Check if the user exists
	checkUser, err := s.DAO.GetUserByEmail(email)
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal server error"})
			return
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "User not found"})
			return
		}
	}

	// Check user status
	if !checkUser.Status {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "message": "Status is Deactivated"})
		return
	}

	// Verify password
	if err := utilities.VerifyPassword(checkUser.Password, password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid email or password"})
		return
	}
	//credentialError := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(password))
	//if credentialError != nil {
	//	log.Println(credentialError.Error())
	//	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	//	ctx.Abort()
	//	return
	//}

	// Generate JWT token
	token, err := Authentication.GenerateJWT(uint(checkUser.ID))
	if err != nil {
		log.Println("err while generating the token", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Error generating JWT token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "token": token, "message": "Login successfully Completed"})
}
func (s *Server) UserProfile(ctx *gin.Context) {
	// Retrieve user ID from the context
	currentUserId := ctx.MustGet("currentUser").(uint)

	// Check if the user ID is valid
	if currentUserId <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "user_id not found"})
		return
	}
	log.Println("currentUserId", currentUserId)
	// Retrieve user details by calling the GetUserByUserID method from the DAO
	result, err := s.DAO.GetUserByUserID(int(currentUserId))
	if err != nil {
		if err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal server error"})
			return
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "User not found"})
			return
		}
	}

	// Return user details in a JSON response
	ctx.JSON(http.StatusOK, gin.H{"data": result, "status": true, "message": "User details fetched successfully"})
}

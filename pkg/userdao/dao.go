package userdao

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"manodarpanNewproject/pkg/models"
)

// DAO defines the interface for user data access operations.
type DAO interface {
	CheckUserID(c *gin.Context, email string) (models.User, error)
	CreateUser(c *gin.Context, user models.User) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByUserID(userID int) (models.User, error)
}

// defaultDAO is the default implementation of DAO.
type defaultDAO struct {
	DBConnector *sqlx.DB
}

func New(db *sqlx.DB) DAO {
	return defaultDAO{
		DBConnector: db,
	}
}

// DBConnector is an interface that provides a DB connection.
type DBConnector interface {
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (Result, error)
}

// Result represents the result of a database operation.
type Result interface {
	RowsAffected() int64
	LastInsertId() int64
}

// CheckUserID checks if a user with the given email exists in the database.
func (d defaultDAO) CheckUserID(c *gin.Context, email string) (models.User, error) {
	user := models.User{}
	err := d.DBConnector.Select(&user, "SELECT * FROM person WHERE first_name=$1", email)
	if err != nil {
		log.Println("error", err)
	}
	fmt.Printf("%#v\n", user)
	return user, err
}
func (d defaultDAO) CreateUser(c *gin.Context, user models.User) (models.User, error) {
	err := d.DBConnector.QueryRowx("INSERT INTO users (first_name,last_name, email,password,phone_no) VALUES ($1, $2,$3,$4,$5) RETURNING id", user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNo).Scan(&user.ID)
	if err != nil {
		log.Println("error", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return models.User{}, err
	}

	return user, nil
}

//func (d defaultDAO) CreateUser(user models.User) (int64, error) {
//	query := `
//		INSERT INTO users (first_name, last_name, email, password, phone_no, status)
//		VALUES ($1,$2,$3,$4,$5,$6)
//	`
//	statement, err := d.DBConnector.Prepare(query)
//	if err != nil {
//		log.Println()
//	}
//	var lastInsertID int64
//	defer statement.Close()
//	err = statement.QueryRow(user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNo, user.Status).Scan(&lastInsertID)
//	if err != nil {
//		log.Println("error while inserting ", err)
//	}
//
//	return lastInsertID, nil
//}

//

// GetUserByEmail retrieves a user from the database by email.
func (d defaultDAO) GetUserByEmail(email string) (user models.User, err error) {
	query := "SELECT * FROM users WHERE email = $1"
	err = d.DBConnector.QueryRowx(query, email).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("messege %s,Error %s,statuscode %s", "Error While Fetching the user details ", err, 500)
			return models.User{}, err
		}

		log.Printf("Error executing query: %v", err)
		return models.User{}, err
	}

	return user, nil
}
func (d defaultDAO) GetUserByUserID(userID int) (user models.User, err error) {
	query := "SELECT * FROM users WHERE id = $1"
	err = d.DBConnector.QueryRowx(query, userID).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("messege %s,Error %s,statuscode %s", "Error While Fetching the user details ", err, 500)
			return models.User{}, err
		}

		log.Printf("Error executing query: %v", err)
		return models.User{}, err
	}

	return user, nil
}

//func (d defaultDAO) GetUserByUserID(userID int) (models.User, error) {
//	var user models.User
//	query := "SELECT * FROM users WHERE email = $1"
//	err := d.DBConnector.Select(&user, query, strconv.Itoa(userID))
//	return user, err
//}

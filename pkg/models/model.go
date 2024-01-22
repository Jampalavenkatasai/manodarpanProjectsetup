package models

import (
	"time"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	PhoneNo   string    `db:"phone_no" json:"phone_no"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Status    bool      `db:"status" json:"status"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	//UUID               string    `db:"uuid" json:"uuid"`
	//SchoolName         string    `db:"school_name" json:"school_name"`
	//Address1           string    `db:"address1" json:"address1"`
	//Address2           string    `db:"address2" json:"address2"`
	//CountryID          int       `db:"country_id" json:"country_id"`
	//StateID            int       `db:"state_id" json:"state_id"`
	//CityID             int       `db:"city_id" json:"city_id"`
	//District           string    `db:"district" json:"district"`
	//Area               string    `db:"area" json:"area"`
	//PostCode           string    `db:"post_code" json:"post_code"`
	//Gender             string    `db:"gender" json:"gender"`
	//Dob                string    `db:"dob" json:"dob"`

	//Active             int       `db:"active" json:"active"`
	//ActiveCode         string    `db:"active_code" json:"active_code"`
	//UserType           int       `db:"user_type" json:"user_type"`
	//ImageName          string    `db:"image_name" json:"image_name"`
	//ClassID            int       `db:"class_id" json:"class_id"`
	//ActiveCodeExpiryAt time.Time `db:"active_code_expiry_at" json:"active_code_expiry_at"`
	//PlanID             JSONMap   `db:"plan_id" json:"plan_id"`
	//PaymentStatus      bool      `db:"payment_status" json:"payment_status"`

	//CreatedBy          int       `db:"created_by" json:"created_by"`
	//UpdatedBy          int       `db:"updated_by" json:"updated_by"`
}

type JSONMap map[string]interface{}

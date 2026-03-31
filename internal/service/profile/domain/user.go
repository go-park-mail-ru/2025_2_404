package user

import (
	// "time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type User struct {
	ID       ID     `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	// HashedPassword string	`json:"password"`
	ImagePath        string `json:"img_path"`
	UserFirstName    string `json:"user_first_name"`
	UserLastName     string `json:"user_second_name"`
	Company          string `json:"company"`
	Phone            string `json:"phone_number"`
	RegistrationDate string `json:"registred_at"`
	AdsCount         int    `json:"ads_count"`
	ProfileType      string `json:"profile_type"`
	CreatedAt        string `json:"created_at"`
}

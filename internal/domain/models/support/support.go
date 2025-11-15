package support

import	(
	modeluser "2025_2_404/internal/domain/models/user"
)

type ID int64

type Support struct {
	ID ID `json:"support_id"`
	UserID modeluser.ID `json:"user_id"`
	Status string `json:"status"`
	Category string `json:"category"`
	Description string `json:"description"`
	ImagePath string `json:"image_path"`
	ContactName string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
}


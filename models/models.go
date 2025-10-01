package models

type BaseUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	BaseUser
	UserName string `json:"user_name"`
}

type Ads struct {
	ID       string `json:"add_id"`
	CreatorID string `json:"creater_id"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

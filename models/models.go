package models

type BaseUser struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type RegisterUser struct {
	BaseUser
	Email string `json:"user_email"`
}

type Ads struct {
	AdID     string `json:"ad_id"`
	CreatorID string `json:"creator_id"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

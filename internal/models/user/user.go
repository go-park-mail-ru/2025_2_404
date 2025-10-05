package user

type BaseUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser struct {
	BaseUser
	UserName string `json:"user_name"`
}

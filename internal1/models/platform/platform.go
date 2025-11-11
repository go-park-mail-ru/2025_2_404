package platform

type ID int

type Platform struct {
	ID  ID `json:"platform_id"`
	Name string `json:"platform_name"`
}
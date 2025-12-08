package slot

import "html/template"

// ID — идентификатор слота.
type ID string

// UserID — идентификатор владельца.
type UserID string

// Slot — доменная модель рекламного слота.
type Slot struct {
	ID              ID
	UserID          UserID
	SlotName        string
	MinCostAdv      int32
	FormatOfBanner  string 
	Status          string 
	BackColor       string 
	TextColor       string 
}

type SlotRenderData struct {
	Title       string
	Description string
	ImageSrc    string
	ImageData   template.URL
	Link        string
	Background  string
	Color       string
}

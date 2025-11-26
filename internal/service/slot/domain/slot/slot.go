package slot

import ()

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
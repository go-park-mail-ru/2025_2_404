// Package slot defines the Slot domain model and related types.
package slot

import "html/template"

type ID string

type UserID string

type Slot struct {
	ID             ID
	UserID         UserID
	SlotName       string
	MinCostAdv     int32
	FormatOfBanner string
	Status         string
	BackColor      string
	TextColor      string
}

type SlotRenderData struct {
	Title       string
	Description string
	ImageSrc    string
	ImageData   template.URL
	Link        string
	Background  string
	Color       string
	Banner      string
	Slot        string
}

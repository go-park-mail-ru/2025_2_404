package models

type Ads struct {
	ID       uint    `json:"add_id"`
	CreatorID string    `json:"creater_id"`
	FilePath string `json:"file_path"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

func (a *Ads) GetID() uint {
	return a.ID
}

func (a *Ads) GetCreatorID() string {
	return a.CreatorID
}

func (a *Ads) GetFilePath() string {
	return a.FilePath
}

func (a *Ads) GetTitle() string {
	return a.Title
}

func (a *Ads) GetText() string {
	return a.Text
}

func (a *Ads) SetID(ID int) {
	a.ID = uint(ID)
}

func (a *Ads) SetCreatorID(CreatorID string) {
	a.CreatorID = CreatorID
}

func (a *Ads) SetFilePath(FilePath string) {
	a.FilePath = FilePath
}

func (a *Ads) SetTitle(Title string) {
	a.Title = Title
}

func (a *Ads) SetText(Text string) {
	a.Text = Text
}

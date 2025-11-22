package storage

type ImageInfo struct {
    ID          string `json:"id"`
    URL         string `json:"url"`
    Filename    string `json:"filename"`
    ContentType string `json:"content_type"`
    Size        int64  `json:"size"`
    UserID      string `json:"user_id"`
    CreatedAt   int64  `json:"created_at"`
    UpdatedAt   int64  `json:"updated_at"`
}
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type VegetaTarget struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Header map[string][]string `json:"header"`
	Body   []byte              `json:"body"`
}

type AdRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	TargetURL string `json:"target_url"`
	Budget    int    `json:"budget"`
	StartAt   string `json:"start_at"`
	EndAt     string `json:"end_at"`
}

func main() {
	token := "Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWE3NTBmNzMtZmRlMS00MjJiLTkyMjYtMzU4YmI4YTZmNGUyIiwiZXhwIjoxNzY2MjQzMTM1LCJpYXQiOjE3NjYxNTY3MzV9.swqRMkq8OPOHQdRsCrw65-tSItcNi_TtAfz6FO9CkcKOKoSas0vT7DSvNQCe9ey0A4oF_8aiq00KqIhR_XCMOA"

	file, err := os.Create("create_ads.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	startTime := time.Now().Format(time.RFC3339)
	endTime := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)

	for i := 0; i < 100000; i++ {
		adPayload := AdRequest{
			Title:     fmt.Sprintf("Sale %d", i),
			Content:   fmt.Sprintf("Super offer number %d", i),
			TargetURL: fmt.Sprintf("https://site.com/promo/%d", i),
			Budget:    rand.Intn(1000) + 100,
			StartAt:   startTime,
			EndAt:     endTime,
		}

		bodyBytes, _ := json.Marshal(adPayload)

		target := VegetaTarget{
			Method: "POST",
			URL:    "http://adnet.website/api/ads",
			Header: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {token},
			},
			Body: bodyBytes,
		}

		// 3. Записываем в файл
		encoder.Encode(target)
	}
	fmt.Println("Done! File 'create_ads.json' generated.")
}

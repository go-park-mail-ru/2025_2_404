// Package convertimage provides utilities for converting images to different formats.
package convertimage

import (
	"encoding/base64"
	"fmt"
	"log"
)

func ConvertImageToBase64(imageData []byte, imageType string) string {
	base64Encoding := base64.StdEncoding.EncodeToString(imageData)
	log.Printf("Converted image to base64: %s", base64Encoding)
	if imageType == "jpg" {
		return fmt.Sprintf("data:image/jpeg;base64,%s", base64Encoding)
	}
	return fmt.Sprintf("data:%s;base64,%s", imageType, base64Encoding)
}

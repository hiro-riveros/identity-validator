package ocr

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func ExtractTextFromImage(imagePath string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(imagePath)
	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text from %s: %w", imagePath, err)
	}

	return text, nil
}

package utils

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func ConvertToSupportedFormat(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".jpg" || ext == ".jpeg" {
		return path, nil
	}

	if ext != ".png" {
		return "", fmt.Errorf("format not supported: %s", ext)
	}

	inputFile, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error while open PNG: %w", err)
	}
	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	if err != nil {
		return "", fmt.Errorf("error decoding PNG: %w", err)
	}

	newPath := strings.TrimSuffix(path, ext) + ".jpg"
	outputFile, err := os.Create(newPath)
	if err != nil {
		return "", fmt.Errorf("error creating JPG: %w", err)
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 95}); err != nil {
		return "", fmt.Errorf("error encoding JPG: %w", err)
	}

	return newPath, nil
}

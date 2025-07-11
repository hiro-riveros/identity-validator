package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-idvalidator/internal/face"
	"go-idvalidator/pkg/validator"
	"log"
)

func main() {
	var (
		frontDNIPath string
		backDNIPath  string
		videoPath    string
		selfiePath   string
	)

	flag.StringVar(&frontDNIPath, "front", "", "Ruta de la imagen delantera del DNI")
	flag.StringVar(&backDNIPath, "back", "", "Ruta de la imagen trasera del DNI")
	flag.StringVar(&videoPath, "video", "", "Ruta del video de rostro")
	flag.StringVar(&selfiePath, "selfie", "", "Ruta de la imagen selfie (opcional)")
	flag.Parse()

	if frontDNIPath == "" || selfiePath == "" {
		log.Fatal("ou must provide at least --front and --selfie")
	}

	recognizer, err := face.NewRecognizer("internal/face/models")
	if err != nil {
		log.Fatalf("Error loading face recognition models: %w", err)
	}
	defer recognizer.Close()

	input := validator.Input{
		FrontDNIPath: frontDNIPath,
		BackDNIPath:  backDNIPath,
		VideoPath:    videoPath,
		SelfiePath:   selfiePath,
	}

	result, err := validator.Analize(input)
	if err != nil {
		log.Fatalf("Error in identity analysis: %v", err)
	}

	fmt.Printf("✅ Resultado:\n")
	fmt.Printf("- Match: %.4f\n", result.MatchPercentage)
	fmt.Printf("- Confidence: %s\n", result.ConfidenceLevel)
	fmt.Printf("- Distance: %.4f\n", result.Distance)
	fmt.Printf("- DNI Text: %+v\n", result.DNIText)
	fmt.Printf("- Reason: %s\n", result.Reason)
	fmt.Printf("- Valid: %v\n", result.Valid)

	jsonOutput, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonOutput))
}

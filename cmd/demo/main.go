package main

import (
	"fmt"
	"go-idvalidator/pkg/validator"
	"log"
)

func main() {
	input := validator.Input{
		FrontDNIPath: "testdata/dni_front_1.png",
		BackDNIPath:  "testdata/dni_back_1.png",
		VideoPath:    "testdata/dni_back_1.mp4",
		SelfiePath:   "testdata/selfie_0.png",
	}

	result, err := validator.Analize(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Resultado:\n")
	fmt.Printf("- Match: %.4f\n", result.MatchPercentage)
	fmt.Printf("- Confidence: %s\n", result.ConfidenceLevel)
	fmt.Printf("- Distance: %.4f\n", result.Distance)
	fmt.Printf("- DNI Text: %+v\n", result.DNIText)
	fmt.Printf("- Reason: %s\n", result.Reason)
	fmt.Printf("- Valid: %v\n", result.Valid)
}

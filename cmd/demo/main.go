package main

import (
	"fmt"
	"go-idvalidator/internal/face"
	validation "go-idvalidator/internal/validator"
	"log"
)

func main() {

	dniPath := "testdata/dni_1.png"
	selfiePath := "testdata/selfie.jpg"

	secret := []byte("AH3TAtrcRAy500VBqXqpwWxf2hdzlpqG")

	recognizer, err := face.NewRecognizer("internal/face/models")
	if err != nil {
		log.Fatalf("error loading models of dlib: %v", err)
	}

	defer recognizer.Close()

	validator := validation.NewValidator(recognizer, secret)
	encrypted, err := validator.ValidateIdentity(dniPath, selfiePath)
	if err != nil {
		log.Fatalf("Identity validation error: %v", err)
	}

	fmt.Println("🔒 Encrypted result:")
	fmt.Println(encrypted)

	result, err := validator.DecryptResult(encrypted)
	if err != nil {
		log.Fatalf("Error decrypting result: %v", err)
	}

	fmt.Println("✅ Decrypted result:")
	fmt.Printf("- distance: %.4f\n", result.Distance)
	fmt.Printf("- Valid coincidence?: %v\n", result.Valid)
}

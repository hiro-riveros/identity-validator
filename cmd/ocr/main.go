package main

import (
	"fmt"
	"go-idvalidator/internal/ocr"
	"go-idvalidator/internal/utils"
)

func main() {

	inputImage := "testdata/dni_front_2.png"
	processedImage, err := ocr.PreprocessImage(inputImage)
	if err != nil {
		fmt.Printf("process error: %v\n", err)
		return
	}
	defer processedImage.Close()

	dni, err := ocr.ExtracDataFromDNI(processedImage)
	if err != nil {
		fmt.Printf("extraction error: %v\n", err)
		return
	}

	fmt.Printf("🆔 DNI Extraído:\n- Name: %s\n- Surname: %s\n- Birthdate: %s\n- Nationality: %s\n- Gender: %s\n- ID Number: %s\n- Document Number: %s\n",
		dni.Name,
		dni.Surname,
		dni.Birthdate,
		dni.Nationality,
		dni.Gender,
		dni.IDNumber,
		dni.DocumentNumber,
	)

	if dni.IDNumber != "" {
		isValid := utils.ValidarRUN(dni.IDNumber)
		fmt.Printf("  -> Invalid RUN: %t\n", isValid)
	}

}

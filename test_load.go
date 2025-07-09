package main

import (
	"fmt"
	"log"

	"github.com/Kagami/go-face"
)

func main() {
	rec, err := face.NewRecognizer("internal/face/models")
	if err != nil {
		log.Fatalf("No se pudo cargar el recognizer: %v", err)
	}
	defer rec.Close()
	fmt.Println("✅ Dlib recognizer cargado correctamente")
}

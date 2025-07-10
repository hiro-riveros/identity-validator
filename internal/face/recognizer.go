package face

import (
	"errors"
	"fmt"

	"github.com/Kagami/go-face"
)

type Recognizer struct {
	recognizer *face.Recognizer
}

func NewRecognizer(modelDir string) (*Recognizer, error) {
	rec, err := face.NewRecognizer(modelDir)
	if err != nil {
		return nil, errors.New("error while loading models: %w")
	}
	return &Recognizer{recognizer: rec}, nil
}

func (rec *Recognizer) Close() {
	rec.recognizer.Close()
}

func (rec *Recognizer) ExtractDescriptor(imagePath string) (face.Descriptor, error) {
	faces, err := rec.recognizer.RecognizeFile(imagePath)
	if err != nil {
		return face.Descriptor{}, fmt.Errorf("error while analyzing image: %w", err)
	}

	if len(faces) == 0 {
		return face.Descriptor{}, fmt.Errorf("no face was detected in %s", imagePath)
	}

	return faces[0].Descriptor, nil
}

func (rec *Recognizer) Compare(d1, d2 face.Descriptor) float64 {
	return face.SquaredEuclideanDistance(d1, d2)
}

func IsEmptyDescriptor(descriptor face.Descriptor) bool {
	for _, vector := range descriptor {
		if vector != 0 {
			return false
		}
	}
	return true
}

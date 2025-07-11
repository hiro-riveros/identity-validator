package validator

import (
	"fmt"
	"go-idvalidator/internal/face"
	"go-idvalidator/internal/validation"
)

var (
	recognizer *face.Recognizer
	secretKey  []byte
	validator  *validation.Validation
)

func init() {
	var err error
	recognizer, err = face.NewRecognizer("internal/face/models")
	if err != nil {
		panic(fmt.Errorf("failed to load recognizer"))
	}

	validator = validation.NewValidation(recognizer, secretKey)
}

func Analize(input Input) (*Result, error) {
	encrypted, err := validator.ValidateIdentity(input.FrontDNIPath, input.SelfiePath)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	response, err := validator.DecryptResult(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decryption error: %w", err)
	}

	return &Result{
		MatchPercentage: response.MatchPercentage,
		ConfidenceLevel: response.ConfidenceLevel,
		Distance:        response.Distance,
		DNIText:         response.DNIText,
		Reason:          response.Reason,
		Valid:           response.Valid,
	}, nil
}

package validation

import (
	"encoding/json"
	"fmt"
	"go-idvalidator/internal/encryption"
	"go-idvalidator/internal/face"
	"go-idvalidator/internal/utils"
)

type Validator struct {
	Face      *face.Recognizer
	SecretKey []byte
}

type Result struct {
	Distance float64 `json:"distance"`
	Valid    bool    `json:"valid"`
}

func NewValidator(face *face.Recognizer, key []byte) *Validator {
	return &Validator{
		Face:      face,
		SecretKey: key,
	}
}

func (validator *Validator) ValidateIdentity(dniFrontPath, selfiePath string) (string, error) {
	dniJPG, err := utils.ConvertToSupportedFormat(dniFrontPath)
	if err != nil {
		return "", fmt.Errorf("DNI convert: %w", err)
	}

	selfieJPG, err := utils.ConvertToSupportedFormat(selfiePath)
	if err != nil {
		return "", fmt.Errorf("DNI convert: %w", err)
	}

	descDNI, err := validator.Face.ExtractDescriptor(dniJPG)
	if err != nil {
		return "", fmt.Errorf("DNI image error: %w", err)
	}

	descSelfie, err := validator.Face.ExtractDescriptor(selfieJPG)
	if err != nil {
		return "", fmt.Errorf("selfie image error: %w", err)
	}

	if face.IsEmptyDescriptor(descDNI) || face.IsEmptyDescriptor(descSelfie) {
		return "", fmt.Errorf("one of the images has no valid face")
	}

	dist := validator.Face.Compare(descDNI, descSelfie)
	result := Result{
		Distance: dist,
		Valid:    dist < 0.6,
	}

	jsonResult, _ := json.Marshal(result)
	encrypted, err := encryption.Encrypt(string(jsonResult), validator.SecretKey)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func (validator *Validator) DecryptResult(encrypted string) (*Result, error) {
	decrypted, err := encryption.Decrypt(encrypted, validator.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("error while decrypting result: %w", err)
	}

	var result Result
	if err := json.Unmarshal([]byte(decrypted), &result); err != nil {
		return nil, fmt.Errorf("error while parsing result: %w", err)
	}

	return &result, nil
}

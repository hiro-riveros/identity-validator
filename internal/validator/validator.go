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
	Distance        float64 `json:"distance"`
	Valid           bool    `json:"valid"`
	MatchPercentage float64 `json:"match_percentage"`
	ConfidenceLevel string  `json:"confidence_level"`
	Reason          string  `json:"reason"`
	DNIText         interface{}
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
	match := distanceToPercentage(dist)
	confidence, valid, reason := translateScore(match)

	result := Result{
		MatchPercentage: match,
		ConfidenceLevel: confidence,
		Distance:        dist,
		Valid:           valid,
		Reason:          reason,
		DNIText: struct {
			Name     string "json:\"name\""
			IDNumber string "json:\"idNumber\""
		}{
			Name:     "",
			IDNumber: "",
		},
	}

	jsonResult, _ := json.MarshalIndent(result, "", " ")
	encrypted, err := encryption.Encrypt(string(jsonResult), validator.SecretKey)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func distanceToPercentage(dist float64) float64 {
	minDist := 0.2
	maxDist := 0.6

	if dist <= minDist {
		return 100.0
	}

	if dist >= maxDist {
		return 0.0
	}

	score := ((maxDist - dist) / (maxDist - minDist) * 100)
	return score
}

func translateScore(score float64) (string, bool, string) {
	switch {
	case score >= 90:
		return "high", true, "High facial match. Features are broadly matched."
	case score >= 75:
		return "medium", true, "Acceptable match. Possible variation in illumination or angle."
	case score >= 60:
		return "low", false, "Low coincidence. Significant differences are detected."
	default:
		return "none", false, "Insufficient facial match"
	}
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

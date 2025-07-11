package validator

type Input struct {
	FrontDNIPath string
	BackDNIPath  string
	VideoPath    string
	SelfiePath   string
}

type Result struct {
	Distance        float64 `json:"distance"`
	Valid           bool    `json:"valid"`
	MatchPercentage float64 `json:"match_percentage"`
	ConfidenceLevel string  `json:"confidence_level"`
	Reason          string  `json:"reason"`
	DNIText         interface{}
}

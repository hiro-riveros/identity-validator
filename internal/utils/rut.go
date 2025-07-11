package utils

import (
	"strconv"
	"strings"
)

// ValidarRUN valida un RUN chileno usando el algoritmo Módulo 11.
func ValidarRUN(run string) bool {
	run = strings.ToUpper(run)
	run = strings.ReplaceAll(run, ".", "")
	run = strings.ReplaceAll(run, "-", "")

	if len(run) < 2 {
		return false
	}

	body := run[:len(run)-1]
	dv := run[len(run)-1:]

	rut, err := strconv.Atoi(body)
	if err != nil {
		return false
	}

	sum := 0
	multiplier := 2
	for ; rut > 0; rut /= 10 {
		sum += (rut % 10) * multiplier
		multiplier++
		if multiplier == 8 {
			multiplier = 2
		}
	}

	calculatedDV := 11 - (sum % 11)

	switch calculatedDV {
	case 11:
		return dv == "0"
	case 10:
		return dv == "K"
	default:
		return dv == strconv.Itoa(calculatedDV)
	}
}

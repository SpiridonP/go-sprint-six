package service

import (
	"strings"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Convert(input string) (string, error) {
	if strings.ContainsAny(input, ".-") {
		return morse.ToText(input), nil
	}
	
	return morse.ToMorse(input), nil
}
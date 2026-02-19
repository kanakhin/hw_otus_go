package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runes := []rune(s)

	output, err := processRunes(runes)

	if err != nil {
		return "", err
	}

	return output, err
}

func repeatRune(r rune, repeatCount int) string {
	output := ""
	for i := 0; i < repeatCount; i++ {
		output += string(r)
	}

	return output
}

func processRunes(runes []rune) (string, error) {
	if len(runes) == 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	if len(runes) == 1 {
		return string(runes[0]), nil
	}

	if string(runes[0]) == "\\" && string(runes[1]) != "\\" && !unicode.IsDigit(runes[1]) {
		return "", ErrInvalidString
	}

	if unicode.IsDigit(runes[1]) && string(runes[0]) == "\\" {
		if !unicode.IsDigit(runes[1]) && string(runes[1]) != "\\" {
			return "", ErrInvalidString
		}

		if len(runes) < 3 {
			output, err := processRunes(runes[2:])

			return string(runes[1]) + output, err
		}

		if !unicode.IsDigit(runes[2]) {
			output, err := processRunes(runes[2:])

			return string(runes[1]) + output, err
		}

		repeatCount, _ := strconv.Atoi(string(runes[2]))
		newOutput, err := processRunes(runes[3:])

		return repeatRune(runes[1], repeatCount) + newOutput, err
	}

	if unicode.IsDigit(runes[1]) {
		repeatCount, _ := strconv.Atoi(string(runes[1]))
		newOutput, err := processRunes(runes[2:])

		return repeatRune(runes[0], repeatCount) + newOutput, err
	}

	if string(runes[0]) == "\\" && string(runes[1]) == "\\" && unicode.IsDigit(runes[2]) {
		repeatCount, _ := strconv.Atoi(string(runes[2]))
		newOutput, err := processRunes(runes[3:])

		return repeatRune(runes[1], repeatCount) + newOutput, err
	}

	if string(runes[0]) == "\\" && string(runes[1]) == "\\" {
		newOutput, err := processRunes(runes[2:])

		return "\\" + newOutput, err
	}

	output := string(runes[0])
	newOutput, err := processRunes(runes[1:])

	return output + newOutput, err
}

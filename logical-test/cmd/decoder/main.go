package main

import (
	"flag"
	"fmt"
	"strconv"
)

// Decoder takes an integer and returns the corresponding letter in the alphabet
// Durring with A=1, B=2, C=3, etc.
// Dombination of intergers is based on sequence of the letters
// 1 -> A
// 12 -> AB, L
// 226 -> BBF, BZ, VF
// 2269 -> BBFI, BZI, VFI, VZ
// If the integer is less than 1 or greater than 26, an error is returned.

type Decoder struct {
	LetterMap map[string]string
}

func (D *Decoder) digitToLetter(digit string) (string, error) {
	fmt.Println("digit being converted to letter:", digit)
	num, err := strconv.Atoi(digit)
	if err != nil {
		return "", err
	}
	if num < 1 || num > 26 {
		return "", fmt.Errorf("number out of range")
	}
	// If the letter is already in the map, return it
	if _, ok := D.LetterMap[digit]; ok {
		return D.LetterMap[digit], nil
	}
	// 'A' is 65 in ASCII, so we add n-1 to it to get the correct letter
	letter := string(rune('A' + num - 1))

	// Add the letter to the map
	D.LetterMap[digit] = letter

	return letter, nil
}

func (D *Decoder) decode(digits string, index int, currentLetter string, result *[]string) error {
	fmt.Println("current letter being decoded:", currentLetter)

	// If index is equal to the length of the digits, we have reached the end of the string
	// And can append the current letter to the result
	if index == len(digits) {
		*result = append(*result, currentLetter)
		return nil
	}

	// Single digit
	letter, err := D.digitToLetter(string(digits[index]))
	if err != nil {
		return err
	}
	// Append letters to combinations
	appendLetters := currentLetter + letter

	// Recursively call decode with the next index and the current letter
	err = D.decode(digits, index+1, appendLetters, result)
	if err != nil {
		return err
	}

	// Check two digits, if valid
	if index < len(digits)-1 {
		// Current digit and the next digit to an integer
		merged := string(digits[index]) + string(digits[index+1])
		digit, err := strconv.Atoi(merged)
		if err != nil {
			return err
		}
		//Need t o convert to an integer and check if the integer is valid
		if digit <= 26 {
			// Use the string value
			letter, err := D.digitToLetter(merged)
			if err != nil {
				return err
			}
			// Append letters again
			appendLetters := currentLetter + letter
			// Recursively call decode with the next index and the current letter to add on to the current letter
			err = D.decode(digits, index+2, appendLetters, result)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (D *Decoder) FindAllCombinations(digits string) (int, error) {
	var result []string
	if digits == "0" {
		// Not valid
		fmt.Println("digits is 0, always return 0")
		result = append(result, "0")
		return len(result), nil
	}
	err := D.decode(digits, 0, "", &result)
	if err != nil {
		return len(result), err
	}
	// Print the result
	fmt.Println("possible combinations:", result)
	return len(result), nil
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func main() {
	// Pass in flag for digits to decode
	digits := flag.String("digits", "1", "digits to decode")
	flag.Parse()

	d := NewDecoder()
	d.LetterMap = make(map[string]string)

	result, err := d.FindAllCombinations(*digits)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("decoded %v to %v\n", *digits, result)

}

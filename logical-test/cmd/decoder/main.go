package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

// Decoder takes an integer and returns the corresponding letter in the alphabet
// Durring with A=1, B=2, C=3, etc.
// Dombination of intergers is based on sequence of the letters
// 1 -> A
// 12 -> AB, L
// 226 -> BBF, BZ, VF
// 2269 -> BBFI, BZI, VFI
// If the integer is less than 1 or greater than 26, an error is returned.

type Decoder struct {
}

func (D *Decoder) digitToLetter(digit string) (string, error) {
	// Convert the digit to an integer
	num, err := strconv.Atoi(digit)
	if err != nil {
		return "", err
	}
	if num < 1 || num > 26 {
		return "", fmt.Errorf("number out of range")
	}
	// Convert the integer to a letter
	return string(rune('A' + num - 1)), nil
}

func (D *Decoder) decode(digits string, index int, currentLetter string, result *[]string) error {
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

func (D *Decoder) FindAllCombinations(digits string) (*int, error) {
	// Make sure the input is valid digits
	_, err := strconv.Atoi(digits)
	if err != nil {
		return nil, fmt.Errorf("digits is not a valid integer")
	}
	// Not valid
	if digits == "0" {
		return nil, fmt.Errorf("digits is 0, always return 0")
	}

	// Add a timer to see how long it takes
	start := time.Now()

	var result []string
	err = D.decode(digits, 0, "", &result)
	if err != nil {
		numberOfCombinations := len(result)
		return &numberOfCombinations, err
	}
	// Print the result
	fmt.Println("possible combinations:", result)
	// Print the time it took
	fmt.Println("time it took:", time.Since(start))
	// Return the number of combinations
	numberOfCombinations := len(result)
	return &numberOfCombinations, nil
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func main() {
	// Pass in flag for digits to decode
	digits := flag.String("digits", "1", "digits to decode")
	flag.Parse()

	d := NewDecoder()

	result, err := d.FindAllCombinations(*digits)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("decoded %v to %v\n", *digits, result)

}

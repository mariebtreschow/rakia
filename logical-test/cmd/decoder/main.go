package main

import (
	"flag"
	"fmt"
	"strconv"
)

// decoder takes an integer and returns the corresponding letter in the alphabet
// curring with A=1, B=2, C=3, etc.
// combination of intergers is based on sequence of the letters
// 1 -> A
// 12 -> AB, L
// 226 -> BBF, BZ, VF
// 2269 -> BBFI, BZI, VFI, VZ
// if the integer is less than 1 or greater than 26, an error is returned.

type Decoder struct {
}

func (D *Decoder) stringNumbersToLetter(s string) (string, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return "", err
	}
	if num < 1 || num > 26 {
		return "", fmt.Errorf("number out of range")
	}
	// 'A' is 65 in ASCII, so we add n-1 to it to get the correct letter
	return string('A' + num - 1), nil
}

func (D *Decoder) decode(digits string, index int, currentLetter string, result *[]string) error {
	fmt.Println("current digit being decoded:", currentLetter)

	// if index is equal to the length of the digits, we have reached the end of the string
	// and can append the current letter to the result
	if index == len(digits) {
		*result = append(*result, currentLetter)
		return nil
	}

	// single digit
	letter, err := D.stringNumbersToLetter(string(digits[index]))
	fmt.Println("1. single letter:", letter)
	if err != nil {
		return err
	}
	// append letters
	appendLetters := currentLetter + letter

	// recursively call decode with the next index and the current letter
	err = D.decode(digits, index+1, appendLetters, result)
	fmt.Println("1. recursively call decode with the appened letters:", appendLetters)
	if err != nil {
		return err
	}

	// two digits, if valid
	if index < len(digits)-1 {
		// current digit and the next digit to an integer
		merged := string(digits[index]) + string(digits[index+1])
		digit, err := strconv.Atoi(merged)
		if err != nil {
			return err
		}
		// need t o convert to an integer and check if the integer is valid
		if digit <= 26 {
			// use the string value
			letter, err := D.stringNumbersToLetter(merged)
			fmt.Println("2. two letters:", merged)
			if err != nil {
				return err
			}
			// append letters
			appendLetters := currentLetter + letter
			// recursively call decode with the next index and the current letter to add on to the current letter
			err = D.decode(digits, index+2, appendLetters, result)
			fmt.Println("2. recursively call decode with the appended letters:", appendLetters)
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
		fmt.Println("digits is 0, always return 0")
		result = append(result, "0")
		return len(result), nil
	}
	err := D.decode(digits, 0, "", &result)
	if err != nil {
		return len(result), err
	}
	fmt.Println("possible combinations:", result)
	return len(result), nil
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func main() {
	digits := flag.String("digits", "1", "digits to decode")
	flag.Parse()

	d := NewDecoder()
	result, err := d.FindAllCombinations(*digits)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Decoded %v to %v\n", *digits, result)

}

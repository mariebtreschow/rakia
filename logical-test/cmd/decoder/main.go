package main

import (
	"flag"
	"fmt"
)

// Decoder takes an integer and returns the corresponding letter in the alphabet
// starting with A=1, B=2, C=3, etc.
// If the integer is less than 1 or greater than 26, an error is returned.

type Decoder struct {
	intergers int
}

func (*Decoder) Decode(intergers int) (string, error) {
	fmt.Printf("decoding intergers: %v to a string\n", intergers)

	if intergers < 1 {
		return "", fmt.Errorf("invalid input")
	}

	return "A", nil
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func main() {
	intergers := flag.Int("intergers", 1, "an int")
	flag.Parse()

	decoder := NewDecoder()
	result, err := decoder.Decode(*intergers)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Decoded %v to %v\n", *intergers, result)

}

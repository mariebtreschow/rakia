# Rakia logical test

Write a function or algorithm that takes a string of digits and returns the number of ways it can be decoded back into its original message.

- Given the input "12", the possible decodings are "AB" and "I", so the output shouldbe 2.
- For the input "226", the possible decodings are "BZ", "VF", and "BBF", making the output 3.
- With the input "O", there are no valid decodings, resulting in an output of 0.

# How to run the program

- Requirements
Go programming language (1.15 or later recommended)

- Installation

Clone the repository to your local machine.
Navigate to the directory containing the code.
Running the Program
To run the decoder use the Go command line:

`go run cmd/decoder/main.go -digits="your_digits_here"`

Replace "your_digits_here" with the digit string you want to decode. For example, to decode "123", you would run:

`go run cmd/decoder/main.go -digits="123"`

## Features

- Decoding Functionality: Converts numeric strings into all possible letter combinations.
- Command-Line Interface: Easy to use command line interface for inputting digit strings.
- Error Handling: Gracefully handles invalid inputs, such as numbers outside the range 1-26.

Usage Examples:
Decode a single number: go run main.go -digits="2"
Decode a sequence of numbers: go run main.go -digits="226"

- Output
The program outputs all possible combinations of the input digits as letters. For example, for the input "226", the output would be the number of combinations followed by the combinations themselves:

`Decoded 226 to 3`

## Tests

To run the test use the Makefile command

`make tests`

## Limitations
The tool currently handles only numbers from 1 to 26 and treats "0" as a special case.
Does not support non-numeric characters.

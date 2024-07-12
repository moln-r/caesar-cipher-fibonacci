package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	supportedInputCharacters       = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,!?@()abcdefghijklmnopqrstuvwxyz ")
	inputFileName                  = "input.txt"
	outputFileName                 = "output.txt"
	whiteSpace                     = ' '
	minAscii                 int32 = 33 // there is 90 characters in total
	maxAscii                 int32 = 122
)

func main() {
	// We'll use all characters from ASCII 33 ('!') to 122 ('z') as potential output characters
	// We'll accept only lower and upper case letters, numbers, and some punctuations ('.', ',', '!', '?', '@', '(', ')')
	// We do not change whitespaces, they will be copied as is

	// If a file contains unsupported characters, we'll log the problem and stop
	// We'll ignore whitespaces at the beginning and end of each line (we trim them before processing)

	// The input must come from input.txt
	// The output will be written to output.txt

	letterToCipherCount, err := validateInputAndCountLetters()
	if err != nil {
		fmt.Println("error validating input:", err)
		return
	}

	fibonacci := generateFibonacci(letterToCipherCount)

	err = encryptAndWriteToFile(&fibonacci)
	if err != nil {
		fmt.Printf("error encrypting and writing to file: %v\n", err)
		return
	}
}

func validateInputAndCountLetters() (int, error) {
	// Init supported character map for validation
	supportedCharactersMap := make(map[rune]bool)
	for _, r := range supportedInputCharacters {
		supportedCharactersMap[r] = true
	}

	file, err := os.Open(inputFileName)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	letterToCipherCount := 0
	preScanner := bufio.NewScanner(file)
	for preScanner.Scan() {
		str := strings.TrimSpace(preScanner.Text())
		for _, r := range str {
			if !supportedCharactersMap[r] && r != whiteSpace {
				return 0, fmt.Errorf("unsupported character: %c\nremove or replace the character before encoding!", r)
			}
		}
		letterToCipherCount += len(str)
	}
	return letterToCipherCount, nil
}

func generateFibonacci(n int) []int {
	var result = make([]int, 0, n)
	result = append(result, 0)
	result = append(result, 1)

	for i := 2; i < n; i++ {
		result = append(result, result[i-1]+result[i-2])
	}

	// because there are 90 characters in total, we'll store mod 90 of the fibonacci numbers
	var i2 = int(maxAscii - minAscii + 1)
	for i := 0; i < len(result); i++ {
		result[i] = result[i] % i2
	}

	return result
}

func encryptAndWriteToFile(fibonacci *[]int) error {
	// Creating the outputFile it if it doesn't exist, override if it does
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error opening or creating outputFile:%w", err)
	}
	defer outputFile.Close()

	inputFileForCypher, err := os.Open(inputFileName)
	if err != nil {
		return fmt.Errorf("Error opening inputFile:%w", err)
	}
	defer inputFileForCypher.Close()
	// Create a new scanner for the inputFile
	scanner := bufio.NewScanner(inputFileForCypher)

	currentLetterIndex := 0 // this will go until letterToCipherCount-1

	// Read the inputFile line by line
	for scanner.Scan() {
		line := scanner.Text()
		encryptedLine := encrypt(&line, &currentLetterIndex, fibonacci)

		// Write data to the outputFile
		_, err = outputFile.WriteString(encryptedLine + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func encrypt(line *string, currentLetterIndex *int, fibonacci *[]int) string {
	var sb strings.Builder
	// go through the line char by char and shift them by the next fibonacci number
	// we'll ignore whitespaces
	for _, r := range *line {
		if r == whiteSpace {
			// we don't change whitespaces
			continue
		}

		// we'll shift the character by the next fibonacci number
		shift := (*fibonacci)[*currentLetterIndex]

		// we'll shift the character by the next fibonacci number
		shiftedRune := r + rune(shift)
		if shiftedRune > maxAscii {
			shiftedRune -= 90
		}
		//fmt.Printf("'%c' shifted by %d is '%c'\n", r, shift, shiftedRune)

		sb.WriteRune(shiftedRune)

		*currentLetterIndex++
	}
	return sb.String()
}

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	enLower = "abcdefghijklmnopqrstuvwxyz"
	enUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ruLower = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	ruUpper = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
)

func shiftRune(r rune, k int) rune {
	findShift := func(alphabet string) rune {
		runes := []rune(alphabet)
		n := len(runes)
		for i, a := range runes {
			if r == a {
				newID := (i + k) % n
				if newID < 0 {
					newID += n
				}
				return runes[newID]
			}
		}
		return r
	}

	switch {
	case (r >= 'a' && r <= 'z'):
		return findShift(enLower)
	case (r >= 'A' && r <= 'Z'):
		return findShift(enUpper)
	case (r >= 'а' && r <= 'я') || r == 'ё':
		return findShift(ruLower)
	case (r >= 'А' && r <= 'Я') || r == 'Ё':
		return findShift(ruUpper)
	default:
		return r
	}
}

func Caesar(str []rune, key int) []rune {
	res := make([]rune, len(str))
	for i, c := range str {
		res[i] = shiftRune(c, key)
	}
	return res
}

func Vigenere(str []rune, key []rune, decrypt bool) []rune {
	res := make([]rune, len(str))
	keyID := 0
	for i, r := range str {
		var shift int
		var found bool

		k := key[keyID%len(key)]

		for i, a := range []rune(enLower) {
			if a == k || []rune(enUpper)[i] == k {
				shift = i
				found = true
				break
			}
		}
		if !found {
			for i, a := range []rune(ruLower) {
				if a == k || []rune(ruUpper)[i] == k {
					shift = i
					found = true
					break
				}
			}
		}

		if decrypt {
			shift = -shift
		}

		newR := shiftRune(r, shift)
		res[i] = newR

		if newR != r {
			keyID++
		}
	}
	return res
}

func main() {
	if len(os.Args) < 6 {
		fmt.Println("Usage: <inFile> <outFile> <enc/dec> <cae/vig> <key>")
		return
	}
	decrypt := false

	fileToRead := os.Args[1]
	fileToWrite := os.Args[2]
	mode := os.Args[3]
	cipherType := os.Args[4]
	keyRaw := os.Args[5]

	var res []rune

	file, err := os.Open(fileToRead)
	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("error reading file: %s\n", err)
		file.Close()
		return
	}
	file.Close()

	input := []rune(string(content))

	if mode == "dec" {
		decrypt = true
	}

	switch cipherType {
	case "cae":
		keyInt, err := strconv.Atoi(keyRaw)
		if err != nil {
			fmt.Printf("key should be positive number for Caesar")
			return
		}
		if decrypt {
			keyInt = -keyInt
		}
		res = Caesar(input, keyInt)
	case "vig":
		if _, err := strconv.Atoi(keyRaw); err == nil {
			fmt.Println("cannot use number for key in Vigenere")
			return
		}
		res = Vigenere(input, []rune(keyRaw), decrypt)
	}

	err = os.WriteFile(fileToWrite, []byte(string(res)), 0644)
	if err != nil {
		fmt.Printf("error writing file: %s\n", err)
	}
	fmt.Printf("Operation completed. Result stored in %s\n\n", fileToWrite)
}

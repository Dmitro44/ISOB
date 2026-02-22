package main

import (
	"fmt"
	"io"
	"os"
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

		shift++

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
	var fileToRead string
	var res []rune

	var choice byte
	for {
		decrypt := false
		for {
			fmt.Println("Choose encrypt or decrypt:")
			fmt.Println("1. Encrypt")
			fmt.Println("2. Decrypt")
			fmt.Println("3. Exit")

			fmt.Scan(&choice)
			switch choice {
			case 1:
				fmt.Println("Enter name of file which contains text to encrypt:")
			case 2:
				fmt.Println("Enter name of file which contains text to decrypt:")
				decrypt = true
			case 3:
				return
			default:
				fmt.Println("Incorrect choice")
				continue
			}
			break
		}
		fmt.Scan(&fileToRead)

		file, err := os.Open(fileToRead)
		if err != nil {
			fmt.Printf("error opening file: %s\n", err)
			continue
		}

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("error reading file: %s\n", err)
			file.Close()
			continue
		}
		file.Close()
		input := []rune(string(content))

		for {
			fmt.Println("\nChoose cipher type: ")
			fmt.Println("1. Caesar")
			fmt.Println("2. Vigenere")
			fmt.Scan(&choice)

			switch choice {
			case 1:
				fmt.Println("Enter key (positive number):")
				var key int
				fmt.Scan(&key)
				if decrypt {
					key = -key
				}
				res = Caesar(input, key)
			case 2:
				fmt.Println("Enter key word:")
				var key string
				fmt.Scan(&key)
				res = Vigenere(input, []rune(key), decrypt)
			default:
				fmt.Println("Incorrect choice")
				continue
			}
			break
		}

		err = os.WriteFile(fileToRead, []byte(string(res)), os.ModeAppend)
		if err != nil {
			fmt.Printf("error writing file: %s\n", err)
		}
		fmt.Printf("Operation completed. Result stored in %s\n", fileToRead)
	}
}

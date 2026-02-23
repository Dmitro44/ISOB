package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(NewModel()).Run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
	// if len(os.Args) < 6 {
	// 	fmt.Println("Usage: <inFile> <outFile> <enc/dec> <cae/vig> <key>")
	// 	return
	// }
	// decrypt := false
	//
	// fileToRead := os.Args[1]
	// fileToWrite := os.Args[2]
	// mode := os.Args[3]
	// cipherType := os.Args[4]
	// keyRaw := os.Args[5]
	//
	// var res []rune
	//
	// file, err := os.Open(fileToRead)
	// if err != nil {
	// 	fmt.Printf("error opening file: %s\n", err)
	// 	return
	// }
	//
	// content, err := io.ReadAll(file)
	// if err != nil {
	// 	fmt.Printf("error reading file: %s\n", err)
	// 	file.Close()
	// 	return
	// }
	// file.Close()
	//
	// input := []rune(string(content))
	//
	// if mode == "dec" {
	// 	decrypt = true
	// }
	//
	// switch cipherType {
	// case "cae":
	// 	keyInt, err := strconv.Atoi(keyRaw)
	// 	if err != nil {
	// 		fmt.Printf("key should be positive number for Caesar")
	// 		return
	// 	}
	// 	if decrypt {
	// 		keyInt = -keyInt
	// 	}
	// 	res = Caesar(input, keyInt)
	// case "vig":
	// 	if _, err := strconv.Atoi(keyRaw); err == nil {
	// 		fmt.Println("cannot use number for key in Vigenere")
	// 		return
	// 	}
	// 	res = Vigenere(input, []rune(keyRaw), decrypt)
	// }
	//
	// err = os.WriteFile(fileToWrite, []byte(string(res)), 0644)
	// if err != nil {
	// 	fmt.Printf("error writing file: %s\n", err)
	// }
	// fmt.Printf("Operation completed. Result stored in %s\n\n", fileToWrite)
}

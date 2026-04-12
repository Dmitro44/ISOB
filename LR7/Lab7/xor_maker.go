package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run xor_maker.go \"string to encrypt\"")
		return
	}

	secret := os.Args[1]
	key := byte(0xAA)

	fmt.Printf("Original:  %s\n", secret)
	fmt.Print("Encrypted: []byte{")

	for i := 0; i < len(secret); i++ {
		encrypted := secret[i] ^ key
		fmt.Printf("0x%02x", encrypted)
		if i < len(secret)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println("}")
}

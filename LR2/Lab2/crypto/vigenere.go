package crypto

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

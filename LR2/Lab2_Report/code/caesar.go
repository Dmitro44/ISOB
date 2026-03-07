package crypto

func Caesar(str []rune, key int) []rune {
	res := make([]rune, len(str))
	for i, c := range str {
		res[i] = shiftRune(c, key)
	}
	return res
}

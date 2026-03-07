package crypto

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

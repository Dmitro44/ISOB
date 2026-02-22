#import "lib/stp2024.typ"
#show: stp2024.template

#include "lab_title.typ"

#stp2024.full_outline()

= Постановка задачи

Целью данной лабораторной работы является получение практических навыков в области классической криптографии путём самостоятельной реализации симметричных шифров замены. Изучение подобных алгоритмов позволяет сформировать базовое понимание принципов защиты информации, оценить стойкость простых шифров и осознать их уязвимости. В качестве объектов реализации выбраны шифр Цезаря и шифр Виженера.

= Выполнение работы

Шифр Цезаря основан на сдвиге каждого символа текста на фиксированное число позиций вдоль алфавита. Числовой ключ $K$ задаёт величину этого сдвига: при шифровании каждая буква заменяется буквой, стоящей на $K$ позиций правее, а при выходе за границу алфавита счёт продолжается с его начала. Дешифрование выполняется обратным сдвигом на то же значение. Поскольку один и тот же ключ применяется ко всему тексту, шифр относится к классу моноалфавитных подстановок и уязвим к частотному анализу. Пример шифрования и дешифрования текста с помощью шифра Цезаря представлен на рисунке @encCae.

#figure(
  image("img/encdec_Caesar.png", width: 75%),
  caption: [Результат шифрования и дешифрования методом Цезаря]
) <encCae>

Шифр Виженера является полиалфавитным обобщением шифра Цезаря. Вместо одного числового ключа используется ключевое слово, каждая буква которого задаёт индивидуальный сдвиг для соответствующей буквы текста. Ключевое слово циклически накладывается на текст: буквы ключевого слова поочерёдно применяются к буквам текста, и после последней буквы ключа счёт начинается заново. Счётчик позиции в ключе продвигается только при обработке литерных символов, благодаря чему знаки препинания и пробелы не нарушают выравнивание ключа. В результате одна и та же буква в разных позициях шифруется по-разному, что существенно усложняет частотный анализ. Пример работы с шифром Виженера показан на рисунке @encVig.

#figure(
  image("img/encdec_Vigenere.png", width: 75%),
  caption: [Результат шифрования и дешифрования методом Виженера]
) <encVig>

Программа реализована в виде консольного приложения. Логика шифрования вынесена в отдельные функции: вспомогательная процедура определяет принадлежность символа к латинскому или русскому алфавиту (строчному или прописному) и выполняет циклический сдвиг с сохранением регистра, нелитерные символы при этом остаются неизменными. На основе этой процедуры построены независимые функции для шифра Цезаря и шифра Виженера. Счётчик позиции в ключе для шифра Виженера продвигается только при обработке литерных символов, что обеспечивает корректное выравнивание ключа при наличии пробелов и знаков препинания. Главный модуль в цикле предлагает пользователю выбрать режим работы (шифрование или дешифрование), ввести имя файла, выбрать алгоритм и указать ключ, результат записывается обратно в тот же файл.

#pagebreak()

#stp2024.heading_unnumbered[Вывод]

В ходе выполнения лабораторной работы реализованы шифры Цезаря и Виженера на языке #emph("Go") с поддержкой русского и латинского алфавитов, сохранением регистра и нелитерных символов. Разработан консольный интерфейс для выбора файла и режима работы. Тестирование подтвердило корректность обоих алгоритмов: последовательное шифрование и дешифрование полностью восстанавливает исходный текст.

#stp2024.appendix(title: [Листинг программного кода], type: [Обязательное])[
  #stp2024.listing[Код программы][
    ```
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
			fmt.Println()

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

		err = os.WriteFile(fileToRead, []byte(string(res)), 0644)
		if err != nil {
			fmt.Printf("error writing file: %s\n", err)
		}
		fmt.Printf("Operation completed. Result stored in %s\n\n", fileToRead)
	}
}
    ```
  ]
]

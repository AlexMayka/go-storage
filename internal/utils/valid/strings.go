package valid

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

var emailWord = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '.',
	'@', '#', '$', '%', '^', '&', '*', '+', ',',
	'1', '2', '3', '4', '5', '6', '7', '8', '9',
	'0',
}

var phoneWord = []rune{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '+',
}

var nameWord = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'а', 'б', 'в', 'г', 'д', 'е', 'ё', 'ж', 'з',
	'и', 'й', 'к', 'л', 'м', 'н', 'о', 'п', 'р',
	'с', 'т', 'у', 'ф', 'х', 'ч', 'ш', 'щ', 'ъ',
	'ы', 'ь', 'э', 'ю', 'я',
}

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
	'е': "e", 'ё': "e", 'ж': "zh", 'з': "z", 'и': "i",
	'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
	'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
	'у': "u", 'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch",
	'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "",
	'э': "e", 'ю': "yu", 'я': "ya",
}

var passwordLetter = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i',
	'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

var passwordChar = []rune{
	'!', '"', '#', '$', '%', '&', '(', ')', '*', '+', '-', '.', '/', ':', ';', '<', '=', '>',
	'?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~', '\'',
}

var passwordNumber = []rune{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
}

func transliterate(s string) string {
	var result strings.Builder
	s = strings.ToLower(s)

	for _, ch := range s {
		if val, ok := translitMap[ch]; ok {
			result.WriteString(val)
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func NormalizationOfName(name string) string {
	name = transliterate(strings.TrimSpace(name))
	re := regexp.MustCompile(`[^\w]+`)
	name = re.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	return name
}

func checkAllowedChar(original string, allow []rune) bool {
	for _, oChar := range original {
		check := false
		for _, aChar := range allow {
			if aChar == oChar {
				check = true
				break
			}

			if aChar != oChar {
				continue
			}
		}

		if !check {
			return false
		}
	}

	return true
}

func CheckEmail(email string) bool {
	if ok := checkAllowedChar(strings.ToLower(email), emailWord); !ok {
		return false
	}
	return true
}

func CheckPhone(phone string) bool {
	if ok := checkAllowedChar(phone, phoneWord); !ok {
		return false
	}
	return true
}

func CheckName(name string) bool {
	if ok := checkAllowedChar(strings.ToLower(name), nameWord); !ok {
		return false
	}
	return true
}

func checkValidPassword(password string) (bool, error) {
	valid := map[string]int{
		"num":     2,
		"char":    1,
		"letter":  3,
		"capital": 1,
	}

	countCh := map[string]int{
		"num":     0,
		"char":    0,
		"letter":  0,
		"capital": 0,
	}

	for _, ch := range password {
		if unicode.IsUpper(ch) {
			countCh["letter"]++
			countCh["capital"]++
			continue
		}

		if slices.Contains(passwordLetter, ch) {
			countCh["letter"]++
			continue
		}

		if slices.Contains(passwordNumber, ch) {
			countCh["num"]++
			continue
		}

		if slices.Contains(passwordChar, ch) {
			countCh["char"]++
			continue
		}

		return false, fmt.Errorf("invalid password")
	}

	for key, count := range countCh {
		if count < valid[key] {
			return false, fmt.Errorf("password must contain at least %d %s", valid[key], key)
		}
	}

	return true, nil
}

func CheckPassword(password string) (bool, error) {
	if len(password) < 6 {
		return false, errors.New("password length must be at least 6 characters")
	}

	if len(password) > 32 {
		return false, errors.New("password length must be less than 32 characters")
	}

	return checkValidPassword(password)
}

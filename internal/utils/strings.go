package utils

import (
	"regexp"
	"strings"
)

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
	'е': "e", 'ё': "e", 'ж': "zh", 'з': "z", 'и': "i",
	'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
	'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
	'у': "u", 'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch",
	'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "",
	'э': "e", 'ю': "yu", 'я': "ya",
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

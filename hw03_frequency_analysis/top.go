package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func clearWord(s string) string {
	runes := []rune(s)

	if len(runes) == 1 && !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
		return ""
	}

	// переводим первый и последний символ строки в нижний регистр
	runes[0] = unicode.ToLower(runes[0])
	runes[len(runes)-1] = unicode.ToLower(runes[len(runes)-1])

	// удаляем первый символ если это знак препинания (ну т.е. не буква и не цифра)
	if !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
		runes = runes[1:]
	}

	// удаляем последний символ если это знак препинания (ну т.е. не буква и не цифра)
	if !unicode.IsLetter(runes[len(runes)-1]) && !unicode.IsDigit(runes[len(runes)-1]) {
		runes = runes[:len(runes)-1]
	}

	return string(runes)
}

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	parts := strings.Fields(text)
	frequency := make(map[string]int)

	for _, word := range parts {
		// "-" не является словом по условию задачи. я расширил это условие, если слово длиной 1
		// и это не буква и не цифра - то я считаю это тоже не словом для определенности
		clearedWord := clearWord(word)
		if clearedWord == "" {
			continue
		}

		frequency[clearedWord]++
	}

	type frequencyStat struct {
		word  string
		count int
	}

	stat := make([]frequencyStat, 0, len(frequency))
	for k, v := range frequency {
		stat = append(stat, frequencyStat{k, v})
	}

	sort.Slice(stat, func(i, j int) bool {
		if stat[i].count == stat[j].count {
			return stat[i].word < stat[j].word
		}

		return stat[i].count > stat[j].count
	})

	top10 := stat
	if len(stat) > 10 {
		top10 = stat[:10]
	}

	output := make([]string, 0, 10)
	for _, v := range top10 {
		output = append(output, v.word)
	}

	return output
}

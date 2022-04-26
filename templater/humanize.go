package templater

import (
	"fmt"
	"math"
	"strings"
)

func ConcatenateWords(words []interface{}, conjunction string) string {
	arrWords := make([]string, len(words))
	for i := range words {
		arrWords[i] = words[i].(string)
	}

	switch len(arrWords) {
	case 0:
		return ""
	case 1:
		return arrWords[0]
	default:
		joinedWords := strings.Join(arrWords[:len(arrWords)-1], ", ")
		lastWord := arrWords[len(arrWords)-1]
		return fmt.Sprintf("%s %s %s", joinedWords, conjunction, lastWord)
	}
}

func Duration(seconds float64) string {
	var humanizedTime string
	intSeconds := int(seconds)

	days := float64(intSeconds) / 60 / 60 / 24
	intSeconds = intSeconds % (60 * 60 * 24)
	hours := float64(intSeconds) / 60 / 60
	intSeconds = intSeconds % (60 * 60)
	minutes := float64(intSeconds) / 60
	intSeconds = intSeconds % 60

	if math.Floor(days) > 0 {
		humanizedTime += Plural(days, "day", "days") + " "
	}
	if math.Floor(hours) > 0 {
		humanizedTime += Plural(hours, "hour", "hours") + " "
	}
	if math.Floor(minutes) > 0 {
		humanizedTime += Plural(minutes, "minute", "minutes") + " "
	}
	if intSeconds > 0 {
		humanizedTime += Plural(float64(intSeconds), "second", "seconds")
	}

	return strings.TrimSpace(humanizedTime)
}

func Plural(quantity float64, singular, plural string) string {
	var suffix = PluralWord(quantity, singular, plural)
	return fmt.Sprintf("%d %s", int(quantity), suffix)
}

func PluralWord(quantity float64, singular, plural string) string {
	if quantity > 1 {
		return plural
	}

	return singular
}

func TruncateQuantity(quantity float64) string {
	if quantity >= 1000 {
		return fmt.Sprintf("%.1fK", quantity/1000)
	}

	return fmt.Sprintf("%d", int(quantity))
}

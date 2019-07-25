package helper

import (
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

func GetStringInBetween(str string, start string, end string) (result *string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)

	if s == -1 || e == -1 {
		return nil
	}

	ret := str[s:e]
	return &ret
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsAlphaNumeric(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]*$")
	return re.MatchString(s)
}

func IsAlphaUnderscore(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z_]*$")
	return re.MatchString(s)
}

func IsStringEqual(text string, characters []string) bool {
	for _, character := range characters {
		if text == character {
			return true
		}
	}
	return false
}

func LastSplit(text string, rune string) []string {
	result := strings.Split(text, rune)

	if len(result) < 2 {
		return []string{
			text,
		}
	}

	return []string{
		strings.Join(result[0:len(result)-1], rune),
		result[len(result)-1],
	}
}
func CalculateRecall(trueSlice interface{}, predictedSlice interface{}) float64 {
	trueReflection := reflect.ValueOf(trueSlice)
	if trueReflection.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	trueValue := make([]interface{}, trueReflection.Len())

	for i := 0; i < trueReflection.Len(); i++ {
		trueValue[i] = trueReflection.Index(i).Interface()
	}

	predictedReflection := reflect.ValueOf(predictedSlice)
	if predictedReflection.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	predictedValue := make([]interface{}, predictedReflection.Len())

	for i := 0; i < predictedReflection.Len(); i++ {
		predictedValue[i] = predictedReflection.Index(i).Interface()
	}

	truePositive := float64(0)
	falseNegative := float64(0)

	for i := 0; i < len(predictedValue)-1; i++ {
		if predictedValue[i] != trueValue[i] {
			falseNegative += 1
		} else {
			truePositive += 1
		}
	}
	return truePositive / (truePositive + falseNegative)
}

func RandFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func RandFloat(min, max float64) float64 {
	res := min + rand.Float64()*(max-min)
	return res
}

func SortByWordValue(wordValue map[string]float64) PairList {
	pl := make(PairList, len(wordValue))
	i := 0
	for k, v := range wordValue {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

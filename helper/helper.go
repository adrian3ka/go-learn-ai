package helper

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode"
)

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

	fmt.Println(predictedValue)
	fmt.Println(trueValue)

	for i := 0; i < len(predictedValue)-1; i++ {
		if predictedValue[i] != trueValue[i] {
			falseNegative += 1
		} else {
			truePositive += 1
		}
	}
	return truePositive / (truePositive + falseNegative)
}

package word_vectorizer

import (
	"regexp"
	"strings"
)

const (
	InvalidType = "Invalid Type"
)

type ArrayWordVectorizer struct {
	lower          bool
	data           map[string]uint64
	labelEncoded   [][2]uint64
	regexReplacers []RegexReplacer
}

type ArrayWordVectorizerConfig struct {
	Lower bool
}

func NewArrayWordVectorizer(vectorizer ArrayWordVectorizerConfig) *ArrayWordVectorizer {
	wv := ArrayWordVectorizer{
		lower: vectorizer.Lower,
		regexReplacers: []RegexReplacer{
			{Pattern: `[^a-zA-Z0-9 ]+`, Replacer: ``},
			{Pattern: `\s+`, Replacer: ` `},
		},
	}

	wv.data = make(map[string]uint64)
	return &wv
}

func (wv *ArrayWordVectorizer) Learn(arrayWord [][2]string) error {
	count := uint64(0)
	for _, pairWord := range arrayWord {
		var word1Vect, word2Vect uint64
		if val, exists := wv.data[pairWord[0]]; !exists {
			wv.data[pairWord[0]] = count
			word1Vect = count
			count++
		} else {
			word1Vect = val
		}
		if val, exists := wv.data[pairWord[1]]; !exists {
			wv.data[pairWord[1]] = count
			word2Vect = count
			count++
		} else {
			word2Vect = val
		}

		tempVectorizedWord := [2]uint64{
			word1Vect,
			word2Vect,
		}

		wv.labelEncoded = append(wv.labelEncoded, tempVectorizedWord)
	}
	return nil
}

func (wv *ArrayWordVectorizer) GetLabelEncodedWords() [][2]uint64 {
	return wv.labelEncoded
}

func (wv *ArrayWordVectorizer) Normalize(document string) (string, error) {
	if wv.lower {
		document = strings.ToLower(document)
	}

	for _, regexReplacer := range wv.regexReplacers {
		reg, err := regexp.Compile(regexReplacer.Pattern)

		if err != nil {
			return "", err
		}

		document = reg.ReplaceAllString(document, regexReplacer.Replacer)
	}
	return document, nil
}

func (wv *ArrayWordVectorizer) GetVectorizedWord() map[string]uint64 {
	return wv.data
}

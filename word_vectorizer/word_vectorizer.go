package word_vectorizer

import (
	"regexp"
	"strings"
)

type RegexReplacer struct {
	Pattern  string
	Replacer string
}

type WordVectorizer struct {
	lower           bool
	data            map[string]uint64 //[word]index
	cleanedCorpuses map[string][]string
	regexReplacers  []RegexReplacer
}

type WordVectorizerConfig struct {
	Lower bool
}

func New(vectorizer WordVectorizerConfig) WordVectorizer {
	cv := WordVectorizer{
		lower: vectorizer.Lower,
		regexReplacers: []RegexReplacer{
			{Pattern: `[^a-zA-Z0-9 ]+`, Replacer: ``},
			{Pattern: `\s+`, Replacer: ` `},
		},
	}

	cv.data = make(map[string]uint64)
	cv.cleanedCorpuses = make(map[string][]string)

	return cv
}

func (cv *WordVectorizer) Learn(corpuses map[string][]string) error {
	for corpusClass, corpus := range corpuses {
		var cleanedCorpus []string
		for _, document := range corpus {
			cleanedDocument := strings.ToLower(document)

			for _, regexReplacer := range cv.regexReplacers {
				reg, err := regexp.Compile(regexReplacer.Pattern)

				if err != nil {
					return err
				}

				cleanedDocument = reg.ReplaceAllString(cleanedDocument, regexReplacer.Replacer)
			}

			tokenizeWords := strings.Split(cleanedDocument, " ")
			for _, word := range tokenizeWords {
				if _, exists := cv.data[word]; !exists {
					cv.data[word] = uint64(len(cv.data))
				}
			}

		}

		if _, exists := cv.cleanedCorpuses[corpusClass]; !exists {
			cv.cleanedCorpuses[corpusClass] = cleanedCorpus
		} else {
			for _, cc := range cleanedCorpus {
				cv.cleanedCorpuses[corpusClass] = append(cv.cleanedCorpuses[corpusClass], cc)
			}
		}
	}
	return nil
}

func (cv *WordVectorizer) GetData() map[string]uint64 {
	return cv.data
}

func (cv *WordVectorizer) GetCleanedCorpus() map[string][]string {
	return cv.cleanedCorpuses
}

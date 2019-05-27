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
	wv := WordVectorizer{
		lower: vectorizer.Lower,
		regexReplacers: []RegexReplacer{
			{Pattern: `[^a-zA-Z0-9 ]+`, Replacer: ``},
			{Pattern: `\s+`, Replacer: ` `},
		},
	}

	wv.data = make(map[string]uint64)
	wv.cleanedCorpuses = make(map[string][]string)

	return wv
}

func (wv *WordVectorizer) Learn(corpuses map[string][]string) error {
	for corpusClass, corpus := range corpuses {
		for _, document := range corpus {

			cleanedDocument, err := wv.Normalize(document)

			if err != nil {
				return err
			}

			tokenizeWords := strings.Split(cleanedDocument, " ")
			for _, word := range tokenizeWords {
				if _, exists := wv.data[word]; !exists {
					wv.data[word] = uint64(len(wv.data))
				}
			}

			wv.cleanedCorpuses[corpusClass] = append(wv.cleanedCorpuses[corpusClass], cleanedDocument)
		}
	}
	return nil
}

func (wv WordVectorizer) Normalize(document string) (string, error) {
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

func (wv WordVectorizer) GetVectorizedWord() map[string]uint64 {
	return wv.data
}

func (wv *WordVectorizer) GetCleanedCorpus() map[string][]string {
	return wv.cleanedCorpuses
}

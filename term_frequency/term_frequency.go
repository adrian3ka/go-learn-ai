package term_frequency

import (
	"strings"
)

type WordVectorizer interface {
	GetVectorizedWord() map[string]uint64
	Normalize(document string) (string, error)
}

type TermFrequency struct {
	data           map[string][][]uint64 //[corpus_name][][]word
	binary         bool
	wordVectorizer WordVectorizer
}

type TermFrequencyConfig struct {
	Binary         bool
	WordVectorizer WordVectorizer
}

func New(config TermFrequencyConfig) TermFrequency {
	tf := TermFrequency{
		binary:         config.Binary,
		wordVectorizer: config.WordVectorizer,
	}

	tf.data = make(map[string][][]uint64)

	return tf
}

func (tf TermFrequency) VectorizedCounter() map[string][][]uint64 {
	return tf.data
}

func (tf TermFrequency) GetDictionary() map[string]uint64 {
	return tf.wordVectorizer.GetVectorizedWord()
}

func (tf TermFrequency) Vectorize(corpusInput []string) ([][]uint64, error) {
	var returnData [][]uint64

	for _, document := range corpusInput {

		cleanedDocument, err := tf.wordVectorizer.Normalize(document)

		if err != nil {
			return nil, err
		}

		slice := tf.countingWord(cleanedDocument)
		returnData = append(returnData, slice)
	}

	return returnData, nil
}
func (tf *TermFrequency) Learn(corpuses map[string][]string) error {

	for corpusClass, corpus := range corpuses {
		for _, document := range corpus {
			slice := tf.countingWord(document)
			tf.data[corpusClass] = append(tf.data[corpusClass], slice)
		}
	}
	return nil
}

func (tf *TermFrequency) countingWord(document string) []uint64 {
	vectorizedWord := tf.wordVectorizer.GetVectorizedWord()

	var slice []uint64
	slice = make([]uint64, len(vectorizedWord))

	tokenizeWords := strings.Split(document, " ")
	for _, word := range tokenizeWords {
		if _, exists := vectorizedWord[word]; exists {
			if tf.binary {
				slice[vectorizedWord[word]] = 1
			} else {
				slice[vectorizedWord[word]] += 1
			}
		}
	}

	return slice
}

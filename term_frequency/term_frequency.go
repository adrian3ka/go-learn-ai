package term_frequency

import (
	"strings"
)

type TermFrequency struct {
	vectorizedWord map[string]uint64
	corpuses       map[string][]string
	data           map[string][][]uint64 //[corpus_name][][]word
	binary         bool
}

type TermFrequencyConfig struct {
	Binary         bool
	VectorizedWord map[string]uint64
}

func New(config TermFrequencyConfig) TermFrequency {
	tf := TermFrequency{
		binary:         config.Binary,
		vectorizedWord: config.VectorizedWord,
	}

	tf.data = make(map[string][][]uint64)

	return tf
}

func (tf *TermFrequency) GetData() map[string][][]uint64 {
	return tf.data
}

func (tf *TermFrequency) Learn(corpuses map[string][]string) error {

	for corpusClass, corpus := range corpuses {
		for _, document := range corpus {
			var slice []uint64
			slice = make([]uint64, len(tf.vectorizedWord))

			tokenizeWords := strings.Split(document, " ")
			for _, word := range tokenizeWords {
				if _, exists := tf.vectorizedWord[word]; exists {
					if tf.binary {
						slice[tf.vectorizedWord[word]] = 1
					} else {
						slice[tf.vectorizedWord[word]] += 1
					}
				}
			}
			tf.data[corpusClass] = append(tf.data[corpusClass], slice)
		}
	}
	return nil
}

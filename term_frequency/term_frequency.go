package term_frequency

import "fmt"

type TermFrequency struct {
	vectorizedWord map[string]uint64
	corpuses       map[string][]string
	data           map[string][][]uint64 //[corpus_name][][]word
	binary         bool
}

type TermFrequencyConfig struct {
	Binary bool
}

func New(config TermFrequencyConfig) TermFrequency {
	tf := TermFrequency{
		binary: config.Binary,
	}

	return tf
}

func (tf *TermFrequency) SetVectorizedWord(input map[string]uint64) {
	tf.vectorizedWord = input
}

func (tf *TermFrequency) Learn(corpuses map[string][]string) error {

	for corpusClass, corpus := range corpuses {
		fmt.Println(corpusClass)
		for _, document := range corpus {
			fmt.Println(document)
		}
	}
	return nil
}

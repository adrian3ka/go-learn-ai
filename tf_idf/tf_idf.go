package tf_idf

import (
	"errors"
	"math"
)

const (
	UnequalDocumentLength = "UnequalDocumentLength"
)

type TermFrequencyInverseDocumentFrequency struct {
	smooth                   bool
	inverseDocumentFrequency []float64
	documentFrequency        []uint64
	data                     map[string][][]uint64
	totalDocument            uint64
	termFrequency            map[string][][]uint64
}

type TermFrequencyInverseDocumentFrequencyConfig struct {
	Smooth        bool
	TermFrequency map[string][][]uint64
}

func New(config TermFrequencyInverseDocumentFrequencyConfig) (TermFrequencyInverseDocumentFrequency, error) {
	tfidf := TermFrequencyInverseDocumentFrequency{
		smooth:        config.Smooth,
		termFrequency: config.TermFrequency,
	}

	tfidf.data = make(map[string][][]uint64)

	//Check Document Length
	var documentLength uint64
	for _, corpuses := range config.TermFrequency {
		for _, corpus := range corpuses {
			tfidf.totalDocument += 1
			if documentLength == 0 {
				documentLength = uint64(len(corpus))
			} else {
				if documentLength != uint64(len(corpus)) {
					return tfidf, errors.New(UnequalDocumentLength)
				}
			}
		}
	}

	tfidf.documentFrequency = make([]uint64, documentLength)
	tfidf.inverseDocumentFrequency = make([]float64, documentLength)

	return tfidf, nil
}

func (tfidf *TermFrequencyInverseDocumentFrequency) Fit() error {
	//Set IDF First
	for _, corpuses := range tfidf.termFrequency {
		for _, corpus := range corpuses {
			for idx, word := range corpus {
				if word != 0 {
					tfidf.documentFrequency[idx] += 1
				}
			}
		}
	}

	for idx, _ := range tfidf.inverseDocumentFrequency {
		numerator := float64(tfidf.totalDocument)
		denominator := float64(tfidf.documentFrequency[idx])

		if tfidf.smooth {
			numerator += 1
			denominator += 1
		}

		tfidf.inverseDocumentFrequency[idx] = math.Log(float64(numerator/denominator)) + 1
	}

	return nil
}

func (tfidf *TermFrequencyInverseDocumentFrequency) GetInverseDocumentFrequency() []float64 {
	return tfidf.inverseDocumentFrequency
}

func (tfidf *TermFrequencyInverseDocumentFrequency) GetDocumentFrequency() []uint64 {
	return tfidf.documentFrequency
}

func (tfidf *TermFrequencyInverseDocumentFrequency) GetData() map[string][][]uint64 {
	return tfidf.data
}

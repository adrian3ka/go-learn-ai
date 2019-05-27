package tf_idf

import (
	"errors"
	"math"
)

const (
	UnequalDocumentLength = "UnequalDocumentLength"

	EuclideanSumSquare = "EuclideanSumSquare"
	EuclideanSum       = "EuclideanSum"
)

type CountVectorizer interface {
	GetDictionary() map[string]uint64
	VectorizedCounter() map[string][][]uint64
	Vectorize([]string) ([][]uint64, error)
}

type TermFrequencyInverseDocumentFrequency struct {
	smooth                   bool
	inverseDocumentFrequency []float64
	documentFrequency        []uint64
	data                     map[string][][]float64
	sumVectorDataPerClass    map[string][]float64
	sumDataPerClass          map[string]float64
	totalDocument            uint64
	normalizerType           string
	normalizer               []float64
	countVectorizer          CountVectorizer
}

type TermFrequencyInverseDocumentFrequencyConfig struct {
	Smooth          bool
	NormalizerType  string
	CountVectorizer CountVectorizer
}

func New(config TermFrequencyInverseDocumentFrequencyConfig) (TermFrequencyInverseDocumentFrequency, error) {
	tfidf := TermFrequencyInverseDocumentFrequency{
		smooth:          config.Smooth,
		normalizerType:  config.NormalizerType,
		countVectorizer: config.CountVectorizer,
	}

	tfidf.data = make(map[string][][]float64)
	tfidf.sumDataPerClass = make(map[string]float64)
	tfidf.sumVectorDataPerClass = make(map[string][]float64)

	//Check Document Length
	var documentLength uint64
	for _, corpuses := range tfidf.countVectorizer.VectorizedCounter() {
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

func (tfidf TermFrequencyInverseDocumentFrequency) Normalize(corpus []uint64) ([]float64, float64) {

	var normalizer = float64(0)
	newTfIdf := make([]float64, len(corpus))
	for idx, word := range corpus {
		newTfIdf[idx] = float64(word) * tfidf.inverseDocumentFrequency[idx]

		if tfidf.normalizerType == EuclideanSumSquare {
			normalizer += math.Pow(newTfIdf[idx], 2)
		} else if tfidf.normalizerType == EuclideanSum {
			normalizer += math.Abs(newTfIdf[idx])
		}
	}

	if normalizer == 0 {
		normalizer = 1
	}

	if tfidf.normalizerType == EuclideanSumSquare {
		normalizer = math.Sqrt(normalizer)
	}

	return newTfIdf, normalizer
}

func (tfidf *TermFrequencyInverseDocumentFrequency) Fit() error {
	//Set IDF First
	for _, corpuses := range tfidf.countVectorizer.VectorizedCounter() {
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

	for corpusClass, corpuses := range tfidf.countVectorizer.VectorizedCounter() {
		sumValue := float64(0)
		sumVectorValue := make([]float64, len(tfidf.documentFrequency))

		for _, corpus := range corpuses {
			newTfIdf, normalizer := tfidf.Normalize(corpus)

			for idx, _ := range corpus {
				newTfIdf[idx] = newTfIdf[idx] / normalizer

				sumVectorValue[idx] += newTfIdf[idx]
				sumValue += newTfIdf[idx]
			}

			tfidf.data[corpusClass] = append(tfidf.data[corpusClass], newTfIdf)
		}
		tfidf.sumVectorDataPerClass[corpusClass] = sumVectorValue
		tfidf.sumDataPerClass[corpusClass] = sumValue
	}

	return nil
}

func (tfidf *TermFrequencyInverseDocumentFrequency) GetInverseDocumentFrequency() []float64 {
	return tfidf.inverseDocumentFrequency
}

func (tfidf *TermFrequencyInverseDocumentFrequency) GetDocumentFrequency() []uint64 {
	return tfidf.documentFrequency
}

func (tfidf TermFrequencyInverseDocumentFrequency) GetSumDataOfClass(class string) float64 {
	if val, exists := tfidf.sumDataPerClass[class]; exists {
		return val
	}

	return 0
}

func (tfidf TermFrequencyInverseDocumentFrequency) GetSumVectorDataOfClass(class string) []float64 {
	if val, exists := tfidf.sumVectorDataPerClass[class]; exists {
		return val
	}

	return nil
}

func (tfidf TermFrequencyInverseDocumentFrequency) GetTrainedData() map[string][][]float64 {
	return tfidf.data
}

func (tfidf TermFrequencyInverseDocumentFrequency) EvaluateInput(input interface{}) ([][]float64, error) {
	var evaluatedInput [][]float64

	convertedInput := input.([]string)

	vectorizedInput, err := tfidf.countVectorizer.Vectorize(convertedInput)

	if err != nil {
		return nil, err
	}

	for _, corpus := range vectorizedInput {
		newTfIdf, normalizer := tfidf.Normalize(corpus)

		for idx, _ := range corpus {
			newTfIdf[idx] = newTfIdf[idx] / normalizer
		}

		evaluatedInput = append(evaluatedInput, newTfIdf)
	}

	return evaluatedInput, nil
}

func (tfidf TermFrequencyInverseDocumentFrequency) GetDictionary() map[string]uint64 {
	return tfidf.countVectorizer.GetDictionary()
}

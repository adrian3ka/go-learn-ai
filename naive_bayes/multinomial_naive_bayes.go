package naive_bayes

import (
	"math"
)

const (
	CONSTANT = 1
)

type EvaluatorInterface interface {
	EvaluateInput(input interface{}) ([][]float64, error)
	GetTrainedData() map[string][][]float64
	GetDictionary() map[string]uint64
	GetSumVectorDataOfClass(class string) []float64
	GetSumDataOfClass(class string) float64
}

type MultinomialNaiveBayesConfig struct {
	Evaluator EvaluatorInterface
}

type MultinomialNaiveBayes struct {
	evaluator EvaluatorInterface
}

func NewMultinomialNaiveBayes(cfg MultinomialNaiveBayesConfig) MultinomialNaiveBayes {
	multinomialNaiveBayes := MultinomialNaiveBayes{
		evaluator: cfg.Evaluator,
	}

	return multinomialNaiveBayes
}

func (nb MultinomialNaiveBayes) Predict(inputs interface{}) ([]string, error) {
	var predicted []string
	probabilities, err := nb.PredictProbability(inputs)

	if err != nil {
		return nil, err
	}

	for _, prob := range probabilities {
		highestProb := float64(0)
		var selectedClass string
		for key, value := range prob {
			if highestProb < value {
				selectedClass = key
				highestProb = value
			}
		}
		predicted = append(predicted, selectedClass)
	}

	return predicted, nil
}

func (nb MultinomialNaiveBayes) PredictProbability(inputs interface{}) ([]map[string]float64, error) {
	evaluatedInputs, err := nb.evaluator.EvaluateInput(inputs)

	if err != nil {
		return nil, err
	}

	var allPrediction []map[string]float64

	for _, evaluatedInput := range evaluatedInputs {
		var predictedClass = make(map[string]float64)
		denominator := float64(0)
		for corpusClass, _ := range nb.evaluator.GetTrainedData() {
			predictedClassValue := float64(1)
			totalValueForClass := nb.evaluator.GetSumDataOfClass(corpusClass)
			dictionaryLength := float64(len(nb.evaluator.GetDictionary()))

			for idx, val := range nb.evaluator.GetSumVectorDataOfClass(corpusClass) {
				predictedWordValue := math.Pow((val+CONSTANT)/(totalValueForClass+dictionaryLength), evaluatedInput[idx])
				predictedClassValue *= predictedWordValue
			}

			denominator += predictedClassValue
			predictedClass[corpusClass] = predictedClassValue
		}

		for corpusClass, predictedValue := range predictedClass {
			predictedClass[corpusClass] = predictedValue / denominator
		}

		allPrediction = append(allPrediction, predictedClass)
	}

	return allPrediction, nil
}

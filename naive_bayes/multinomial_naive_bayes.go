package naive_bayes

import "fmt"

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

func (nb MultinomialNaiveBayes) Predict(input interface{}) (map[string]float64, error) {
	var predictedValue = make(map[string]float64)

	evaluatedInput, err := nb.evaluator.EvaluateInput(input)

	if err != nil {
		return nil, err
	}

	fmt.Println(evaluatedInput)
	for corpusClass, corpuses := range nb.evaluator.GetTrainedData() {
		fmt.Println(corpusClass, corpuses)

		fmt.Println(nb.evaluator.GetSumDataOfClass(corpusClass), nb.evaluator.GetSumVectorDataOfClass(corpusClass))
	}

	return predictedValue, nil
}

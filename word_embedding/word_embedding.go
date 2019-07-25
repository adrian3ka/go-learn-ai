package word_embedding

import (
	"errors"
	"fmt"
	"github.com/adrian3ka/go-learn-ai/helper"
	"math/rand"
	"time"
)

type LossFunctionType string
type ActivationFunctionType string
type OptimizerType string

const (
	OneHotEncodedNil = "One Hot Encoded Data is Nil"
	InvalidDataLearn = "Invalid Data Learn"

	CrossEntropy             LossFunctionType       = "cross_entropy"
	Softmax                  ActivationFunctionType = "softmax"
	GradientDescentOptimizer OptimizerType          = "gradient_descent_optimizer"
)

type Layer struct {
	Value  [][]float64
	Weight [][]float64
	Bias   float64
}

type WordEmbedding struct {
	dimension          uint64
	hiddenLayer        Layer
	outputLayer        Layer
	lossFunction       LossFunctionType
	activationFunction ActivationFunctionType
	optimizer          OptimizerType
	oneHotEncodedData  OneHotEncoderInterface
	learningRate       float64
}

type OneHotEncoderInterface interface {
	GetEncodedData() [][2][]uint64
}

type WordEmbeddingConfig struct {
	Dimension          uint64
	OneHotEncodedData  OneHotEncoderInterface
	LossFunction       LossFunctionType
	ActivationFunction ActivationFunctionType
	Optimizer          OptimizerType
	LearningRate       float64
}

func (we *WordEmbedding) Learn(source, target [][]float64) error {
	if len(source) != len(target) {
		return errors.New(InvalidDataLearn)
	}

	for idx := range source {
		var tempSource, tempTarget [][]float64
		tempSource = append(tempSource, source[idx])
		tempTarget = append(tempTarget, target[idx])
		tempMatrix1, err := helper.MatrixMultiplication(tempSource, we.hiddenLayer.Weight)

		if err != nil {
			return err
		}

		we.hiddenLayer.Value = helper.MatrixAdditionWithNumber(tempMatrix1, we.hiddenLayer.Bias)

		tempMatrix2, err := helper.MatrixMultiplication(we.hiddenLayer.Value, we.outputLayer.Weight)

		if err != nil {
			return err
		}

		tempMatrix2 = helper.MatrixAdditionWithNumber(tempMatrix2, we.outputLayer.Bias)

		if we.activationFunction == Softmax {
			tempMatrix2[0] = helper.Softmax(tempMatrix2[0])
		}
		fmt.Println(tempMatrix2)
	}

	return nil
}

func NewWordEmbedding(config WordEmbeddingConfig) (*WordEmbedding, error) {

	if config.LearningRate == 0 {
		config.LearningRate = 1
	}
	rand.Seed(time.Now().UTC().UnixNano())
	if config.OneHotEncodedData == nil {
		return nil, errors.New(OneHotEncodedNil)
	}

	if len(config.OneHotEncodedData.GetEncodedData()) == 0 {
		return nil, errors.New(OneHotEncodedNil)
	}

	hiddenLayer := Layer{
		Bias: helper.RandFloat(0, 1),
	}

	for i := 0; i < len(config.OneHotEncodedData.GetEncodedData()[0][0]); i++ {
		tempRow := helper.RandFloats(0, 1, int(config.Dimension))
		hiddenLayer.Weight = append(hiddenLayer.Weight, tempRow)
	}

	outputLayer := Layer{
		Bias: helper.RandFloat(0, 1),
	}

	for i := 0; i < int(config.Dimension); i++ {
		tempRow := helper.RandFloats(0, 1, len(config.OneHotEncodedData.GetEncodedData()[0][0]))
		outputLayer.Weight = append(outputLayer.Weight, tempRow)
	}

	we := WordEmbedding{
		learningRate:       config.LearningRate,
		outputLayer:        outputLayer,
		hiddenLayer:        hiddenLayer,
		oneHotEncodedData:  config.OneHotEncodedData,
		dimension:          config.Dimension,
		optimizer:          config.Optimizer,
		lossFunction:       config.LossFunction,
		activationFunction: config.ActivationFunction,
	}

	var source, target [][]float64
	for _, learnData := range we.oneHotEncodedData.GetEncodedData() {
		if len(learnData[0]) != len(learnData[1]) {
			return nil, errors.New(InvalidDataLearn)
		}
		var tempSource, tempTarget []float64
		for idx := range learnData[0] {
			tempSource = append(tempSource, float64(learnData[0][idx]))
			tempTarget = append(tempTarget, float64(learnData[1][idx]))
		}
		source = append(source, tempSource)
		target = append(target, tempTarget)
	}

	err := we.Learn(source, target)

	if err != nil {
		return nil, err
	}

	return &we, nil
}

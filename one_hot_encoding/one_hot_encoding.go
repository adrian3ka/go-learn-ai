package one_hot_encoding

type ArrayWordVectorizerInterface interface {
	GetVectorizedWord() map[string]uint64
	GetLabelEncodedWords() [][2]uint64
}

type OneHotEncoder struct {
	oneHotEncodedData [][2][]uint64
	labelEncodedData  ArrayWordVectorizerInterface
}

type OneHotEncoderConfig struct {
	LabelEncodedData ArrayWordVectorizerInterface
}

func NewOneHotEncoder(config OneHotEncoderConfig) *OneHotEncoder {
	o := OneHotEncoder{
		labelEncodedData: config.LabelEncodedData,
	}

	oneHotEncoderArrayLength := len(o.labelEncodedData.GetVectorizedWord())

	var skeletonArray []uint64
	for i := 0; i < oneHotEncoderArrayLength; i++ {
		skeletonArray = append(skeletonArray, 0)
	}

	for _, labeledWords := range o.labelEncodedData.GetLabelEncodedWords() {
		var tempArray, tempArray2 []uint64

		for i := 0; i < oneHotEncoderArrayLength; i++ {
			tempArray = append(tempArray, 0)
			tempArray2 = append(tempArray2, 0)
		}

		tempArray[labeledWords[0]] = 1
		tempArray2[labeledWords[1]] = 1

		tempPairArray := [2][]uint64{
			tempArray,
			tempArray2,
		}
		o.oneHotEncodedData = append(o.oneHotEncodedData, tempPairArray)
	}

	return &o
}

func (o *OneHotEncoder) GetEncodedData() [][2][]uint64 {
	return o.oneHotEncodedData
}

func (o *OneHotEncoder) Encode([][2]string) [][2]float64 {
	return nil
}

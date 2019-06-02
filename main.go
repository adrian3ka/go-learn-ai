package main

import (
	"fmt"
	"github.com/adrian/go-learn-ai/helper"
	"github.com/adrian/go-learn-ai/naive_bayes"
	"github.com/adrian/go-learn-ai/tagger"
	"github.com/adrian/go-learn-ai/term_frequency"
	"github.com/adrian/go-learn-ai/tf_idf"
	"github.com/adrian/go-learn-ai/word_vectorizer"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println("============================= Classifier =====================================")
	wordVectorizer := word_vectorizer.New(word_vectorizer.WordVectorizerConfig{
		Lower: true,
	})

	var corpuses map[string][]string

	corpuses = make(map[string][]string)

	corpuses["pulsa"] = []string{
		"Saya mau beli pulsa dong. Jual voucher gak bang?. Mau isi pulsa dong.",
		"jual pulsa gak ya?",
		"kamu jual voucher ga?",
		"mau isi paket data bisa?",
		"mau isi pulsa bisa ga ya?",
	}

	corpuses["tiket"] = []string{
		"kamu jual tiket pesawat ga?",
		"disini jual tiket ga ya?",
		"bisa beli tiket    kereta?",
		"jual tiket apa  ya?",
	}
	err := wordVectorizer.Learn(corpuses)

	if err != nil {
		panic(err)
	}

	termFrequency := term_frequency.New(term_frequency.TermFrequencyConfig{
		Binary:         false,
		WordVectorizer: wordVectorizer,
	})

	err = termFrequency.Learn(wordVectorizer.GetCleanedCorpus())

	if err != nil {
		panic(err)
	}

	tfIdf, err := tf_idf.New(tf_idf.TermFrequencyInverseDocumentFrequencyConfig{
		Smooth:          true,
		NormalizerType:  tf_idf.EuclideanSumSquare,
		CountVectorizer: termFrequency,
	})

	if err != nil {
		panic(err)
	}

	err = tfIdf.Fit()

	if err != nil {
		panic(err)
	}

	multinomialNB := naive_bayes.NewMultinomialNaiveBayes(naive_bayes.MultinomialNaiveBayesConfig{
		Evaluator: tfIdf,
	})

	predicted, err := multinomialNB.Predict([]string{
		"mAu belI tiket kEreta doNg",
		"jual pulsa ga ya?",
		"mau beli tiket kereta pake pulsa bisa ga ya?",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(predicted)

	fmt.Println("============================== POS Tagger ====================================")
	file, err := ioutil.ReadFile("tagged_corpus/Indonesian.txt")
	if err != nil {
		log.Fatal(err)
	}

	defaultTag := "NN"
	allTuple := tagger.StringToTuple(tagger.StringToTupleInput{
		Text:    string(file),
		Lower:   false,
		Default: &defaultTag,
	})

	border := len(allTuple.Tuple) * 99 / 100
	trainTuple := allTuple.Tuple[0:border]
	testTuple := allTuple.Tuple[border:len(allTuple.Tuple)]

	var testTaggedWord [][2]string
	testSentence := ""

	for _, sentence := range testTuple {
		for _, word := range sentence {
			testTaggedWord = append(testTaggedWord, word)
			testSentence += word[0] + " "
		}
	}

	defaultTagger := tagger.NewDefaultTagger(tagger.DefaultTaggerConfig{
		DefaultTag: "nn",
	})

	err = defaultTagger.Learn(trainTuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err := defaultTagger.Predict(testSentence)

	if err != nil {
		panic(err)
	}

	fmt.Println("Recall Of Default Tagger Only >> ", helper.CalculateRecall(testTaggedWord, predictedValue))

	unigramTagger := tagger.NewUnigramTagger(tagger.UnigramTaggerConfig{
		BackoffTagger: defaultTagger,
	})

	err = unigramTagger.Learn(trainTuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = unigramTagger.Predict(testSentence)

	if err != nil {
		panic(err)
	}

	fmt.Println("Recall Of Unigram Tagger With Backoff >> ", helper.CalculateRecall(testTaggedWord, predictedValue))

	bigramTagger := tagger.NewNGramTagger(tagger.NGramTaggerConfig{
		BackoffTagger: unigramTagger,
		N:             2,
	})

	err = bigramTagger.Learn(trainTuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = bigramTagger.Predict(testSentence)

	if err != nil {
		panic(err)
	}

	fmt.Println("Recall Of Bigram Tagger With Backoff >> ", helper.CalculateRecall(testTaggedWord, predictedValue))

	trigramTagger := tagger.NewNGramTagger(tagger.NGramTaggerConfig{
		BackoffTagger: bigramTagger,
		N:             3,
	})

	err = trigramTagger.Learn(trainTuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = trigramTagger.Predict(testSentence)

	if err != nil {
		panic(err)
	}

	fmt.Println("Recall Of Trigram Tagger With Backoff >> ", helper.CalculateRecall(testTaggedWord, predictedValue))

	regexTagger := tagger.NewRegexTagger(tagger.RegexTaggerConfig{
		Patterns:      tagger.DefaultSimpleIndonesianRegexTagger,
		BackoffTagger: trigramTagger,
	})

	err = regexTagger.Learn(trainTuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = regexTagger.Predict(testSentence)

	if err != nil {
		panic(err)
	}

	fmt.Println("Recall Of Regex Tagger With Backoff >> ", helper.CalculateRecall(testTaggedWord, predictedValue))
}

package main

import (
	"fmt"
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

	fmt.Println(wordVectorizer.GetVectorizedWord())

	termFrequency := term_frequency.New(term_frequency.TermFrequencyConfig{
		Binary:         false,
		WordVectorizer: wordVectorizer,
	})

	err = termFrequency.Learn(wordVectorizer.GetCleanedCorpus())

	if err != nil {
		panic(err)
	}

	fmt.Println(termFrequency.VectorizedCounter())

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

	fmt.Println(tfIdf.GetInverseDocumentFrequency())
	fmt.Println(tfIdf.GetDocumentFrequency())
	fmt.Println(tfIdf.GetTrainedData())

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

	fmt.Println("=============================== POS Tagger =====================================")
	file, err := ioutil.ReadFile("tagged_corpus/Indonesian.txt")
	if err != nil {
		log.Fatal(err)
	}
	tuple := tagger.StringToTuple(tagger.StringToTupleInput{
		Text:  string(file),
		Lower: true,
	})

	text := "Halo nama saya Adrian. Saya baik dan tampan. Saya berumur 10 tahun."
	defaultTagger := tagger.NewDefaultTagger(tagger.DefaultTaggerConfig{
		DefaultTag: "nn",
	})

	err = defaultTagger.Learn(tuple.Tuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err := defaultTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	fmt.Println(predictedValue)

	regexTagger := tagger.NewRegexTagger(tagger.RegexTaggerConfig{
		Patterns:      tagger.DefaultIndonesianRegexTagger,
		BackoffTagger: defaultTagger,
	})

	err = regexTagger.Learn(tuple.Tuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = regexTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	fmt.Println(predictedValue)

	unigramTagger := tagger.NewUnigramTagger(tagger.UnigramTaggerConfig{
		BackoffTagger: regexTagger,
	})

	err = unigramTagger.Learn(tuple.Tuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = unigramTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	fmt.Println(predictedValue)

	bigramTagger := tagger.NewNGramTagger(tagger.NGramTaggerConfig{
		BackoffTagger: unigramTagger,
		N:             2,
	})

	err = bigramTagger.Learn(tuple.Tuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = bigramTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	fmt.Println(predictedValue)

	trigramTagger := tagger.NewNGramTagger(tagger.NGramTaggerConfig{
		BackoffTagger: bigramTagger,
		N:             3,
	})

	err = trigramTagger.Learn(tuple.Tuple)

	if err != nil {
		panic(err)
	}

	predictedValue, err = trigramTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	fmt.Println(predictedValue)
}

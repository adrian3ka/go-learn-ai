package main

import (
	"fmt"
	"github.com/adrian3ka/go-learn-ai/grammar_parser"
	"github.com/adrian3ka/go-learn-ai/helper"
	"github.com/adrian3ka/go-learn-ai/naive_bayes"
	nfa2 "github.com/adrian3ka/go-learn-ai/nfa"
	"github.com/adrian3ka/go-learn-ai/tagger"
	"github.com/adrian3ka/go-learn-ai/term_frequency"
	"github.com/adrian3ka/go-learn-ai/tf_idf"
	"github.com/adrian3ka/go-learn-ai/word_vectorizer"
	"io/ioutil"
	"log"
	"os"
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

	corpuses["saldo"] = []string{
		"halo aku mau isi saldo dong",
		"eh mau topup dong bisa ga?",
		"mau nambah saldo dong bisa gak",
		"tolong bantu isi saldo dong 50 ribu",
	}

	corpuses["hotel"] = []string{
		"ada kamar kosong ga ya di tempat kamu?",
		"eh mau sewa 1 kamar dong",
		"eh  mau sewa ruang meeting dong",
		"eh fasilitas di hotel kamu apa aja ya?",
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

	dataTest := []string{
		"mAu belI tiket kEreta doNg",
		"jual pulsa ga ya?",
		"mau beli tiket kereta pake pulsa bisa ga ya?",
		"mau topup isi wallet dong",
		"saya mau beli pulsa dong",
		"Ada ruang meeting ga di hotel ini ?",
	}

	predicted, err := multinomialNB.Predict(dataTest)

	if err != nil {
		panic(err)
	}

	for idx, p := range predicted {
		fmt.Println(dataTest[idx], " >> ", p)
	}

	fmt.Println("============================== POS Tagger ====================================")
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.ReadFile(dir + "/go-learn-ai/tagged_corpus/Indonesian.txt")
	if err != nil {
		log.Fatal(err)
	}

	defaultTag := "NN"
	allTuple := tagger.StringToTuple(tagger.StringToTupleInput{
		Text:    string(file),
		Lower:   false,
		Default: &defaultTag,
	})

	border := len(allTuple.Tuple) * 90 / 100
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

	regexTagger := tagger.NewRegexTagger(tagger.RegexTaggerConfig{
		Patterns:      tagger.DefaultSimpleIndonesianRegexTagger,
		BackoffTagger: defaultTagger,
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

	unigramTagger := tagger.NewUnigramTagger(tagger.UnigramTaggerConfig{
		BackoffTagger: regexTagger,
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

	fmt.Println("=================================== NFA ======================================")
	nfa, state0, err := nfa2.NewNFA("State 0", false)

	if err != nil {
		panic(err)
	}

	state1, err := nfa.AddState(&nfa2.State{
		Name: "State 1",
	}, false)

	if err != nil {
		panic(err)
	}

	state2, err := nfa.AddState(&nfa2.State{
		Name: "State 2",
	}, true)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	err = nfa.AddTransition(state0.Index, "a", *state1, *state2)

	if err != nil {
		panic(err)
	}

	err = nfa.AddTransition(state1.Index, "b", *state0, *state2)

	if err != nil {
		panic(err)
	}

	var inputs []string
	fmt.Println("Input a")
	inputs = append(inputs, "a")

	fmt.Println("Input b")

	inputs = append(inputs, "b")

	nfa.PrintTransitionTable()

	fmt.Println("If input a, b will go to final?", nfa.VerifyInputs(inputs))
	fmt.Println("If input a, b will go to final?", nfa.VerifyInputs(inputs))

	fmt.Println("========================= Information Extraction =============================")

	text := "Menteri Perhubungan Ignasius Jonan dan Jaksa Agung Prasetyo , menandatangani MoU tentang kordinasi dalam pelaksanaan tugas dan fungsi"

	predictedValue, err = bigramTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	fmt.Println(predictedValue)

	gp, err := grammar_parser.NewRegexpParser(grammar_parser.RegexpParserConfig{
		Grammar: [][2]string{
			{"NP", "{<NN>+}"}, //Chunking
		},
	})

	if err != nil {
		panic(err)
	}

	parsedGrammar, err := gp.Parse(predictedValue)

	if err != nil {
		panic(err)
	}

	for _, x := range parsedGrammar {
		if x.GeneralTag != nil {
			fmt.Println(*x.GeneralTag, "-> ", x.Words)
		} else {
			fmt.Println("-> ", x.Words)
		}
	}

	fmt.Println("-----------------------------------------------------------------------")
	text = "Nama saya Adrian Eka Sanjaya ."

	predictedValue, err = bigramTagger.Predict(text)

	if err != nil {
		panic(err)
	}

	parsedGrammar, err = gp.Parse(predictedValue)

	if err != nil {
		panic(err)
	}

	for _, x := range parsedGrammar {
		if x.GeneralTag != nil {
			fmt.Println(*x.GeneralTag, "-> ", x.Words)
		} else {
			fmt.Println("-> ", x.Words)
		}
	}
	fmt.Println("-----------------------------------------------------------------------")
}

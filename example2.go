package main

import (
	"fmt"
	"github.com/adrian3ka/go-learn-ai/grammar_parser"
)

func main() {
	fmt.Println("==================================CHUNKING===========================================")
	taggedSentence := [][2]string{
		{"the", "DT"},
		{"little", "JJ"},
		{"yellow", "JJ"},
		{"dog", "NN"},
		{"barked", "VBD"},
		{"at", "IN"},
		{"the", "DT"},
		{"cat", "NN"},
	}

	gp, err := grammar_parser.NewRegexpParser(grammar_parser.RegexpParserConfig{
		Grammar: [][2]string{
			{"NP", "{<DT>?<JJ>*<NN>}"}, //Chunking
		},
	})

	if err != nil {
		panic(err)
	}

	parsedGrammar, err := gp.Parse(taggedSentence)

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

	taggedSentence2 := [][2]string{
		{"Rapunzel", "NNP"},
		{"let", "VBD"},
		{"down", "RP"},
		{"her", "PRP"},
		{"long", "JJ"},
		{"golden", "JJ"},
		{"hair", "NN"},
	}

	gp2, err := grammar_parser.NewRegexpParser(grammar_parser.RegexpParserConfig{
		Grammar: [][2]string{
			{"NP", "{<DT|PRP>?<JJ>*<NN>}"}, //Chunking
			{"NP", "{<NNP>+}"},             //Chunking
		},
	})

	if err != nil {
		panic(err)
	}

	parsedGrammar, err = gp2.Parse(taggedSentence2)

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

	fmt.Println("==================================CHINKING===========================================")

	gp3, err := grammar_parser.NewRegexpParser(grammar_parser.RegexpParserConfig{
		Grammar: [][2]string{
			{"NP", "}<VBD><RP>{"}, //Chinking
		},
	})

	if err != nil {
		panic(err)
	}

	parsedGrammar, err = gp3.Parse(taggedSentence2)

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
}

package main

import (
	"fmt"
	"github.com/adrian/go-learn-ai/grammar_parser"
)

func main() {
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
}

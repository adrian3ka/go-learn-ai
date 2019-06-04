package grammar_parser

import (
	"errors"
	"fmt"
	"github.com/adrian/go-learn-ai/helper"
	"github.com/adrian/go-learn-ai/nfa"
	"strings"
)

//AVAIABLE SYMBOL
//- CHINCKING : } {
//- CHUNGKING : { }
//- One Or More : +

const (
	CannotCreateNFAFromGrammar = "Cannot Create NFA From Grammar"
	InvalidGrammar             = "Invalid Grammar"
	NilState                   = "Nil State"
	OneOrMore                  = "+"
	OpeningTag                 = "<"
	ClosingTag                 = ">"
)

type BasicParser interface {
	Parse([][2]string) error
}

type NfaGrammar struct {
	Nfa    nfa.NFA
	Target string
}

type RegexpParser struct {
	nfaGrammar []NfaGrammar
}

type RegexpParserConfig struct {
	Grammar [][2]string
}

func handleOneOrMore(nfaData *nfa.NFA, tag string, isFinal bool) error {

	state1, err := nfaData.AddState(&nfa.State{
		Name: tag,
	}, false)

	if err != nil {
		return err
	}

	state2, err := nfaData.AddState(&nfa.State{
		Name: tag,
	}, isFinal)

	if err != nil {
		return err
	}

	if state1 == nil || state2 == nil {
		return errors.New(NilState)
	}

	nfaData.AddTransition(state1.Index, tag, *state2)
	return nil
}

func convertGrammarToNfa(grammar string) (*nfa.NFA, error) {
	fmt.Println(grammar)

	var newNFA *nfa.NFA
	var err error

	if string(grammar[0]) == "{" && string(grammar[len(grammar)-1]) == "}" {
		grammar = strings.Replace(grammar, "{", "", 1)
		grammar = strings.Replace(grammar, "}", "", 1)
		fmt.Println(grammar)
		var nextTag *string
		isFinal := false
		for {
			tag := helper.GetStringInBetween(grammar, OpeningTag, ClosingTag)

			if tag == nil {
				break
			}

			processedTag := OpeningTag + *tag + ClosingTag

			fmt.Println(processedTag)

			grammar = strings.Replace(grammar, processedTag, "", 1)

			nextTag = helper.GetStringInBetween(grammar, "<", ">")

			if nextTag == nil {
				fmt.Println("isFinal >> ", isFinal)
				isFinal = true
			}

			var state *nfa.State

			newNFA, state, err = nfa.NewNFA(grammar, isFinal)

			if err != nil {
				return nil, err
			}

			fmt.Println(*state)

			fmt.Println("Grammar >> ", grammar)
		}
	} else {
		return nil, errors.New(InvalidGrammar)
	}
	return newNFA, nil
}

func NewRegexpParser(config RegexpParserConfig) (*RegexpParser, error) {
	var nfas []NfaGrammar

	for _, g := range config.Grammar {
		newNfa, err := convertGrammarToNfa(g[1])

		if err != nil {
			return nil, err
		}

		if newNfa == nil {
			return nil, errors.New(CannotCreateNFAFromGrammar)
		}

		nfas = append(nfas, NfaGrammar{
			Nfa:    *newNfa,
			Target: g[0],
		})
	}
	return &RegexpParser{
		nfaGrammar: nfas,
	}, nil
}

func (rp *RegexpParser) Parse(input [][2]string) error {
	fmt.Println(input)
	return nil
}

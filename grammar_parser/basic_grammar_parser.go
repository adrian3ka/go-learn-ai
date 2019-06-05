package grammar_parser

import (
	"errors"
	"github.com/adrian/go-learn-ai/helper"
	"github.com/adrian/go-learn-ai/nfa"
	"strings"
)

//AVAIABLE SYMBOL
//- CHINCKING : } {
//- CHUNGKING : { }

const (
	CannotCreateNFAFromGrammar = "Cannot Create NFA From Grammar"
	InvalidGrammar             = "Invalid Grammar"
	NilState                   = "Nil State"
	OneOrMore                  = "+"
	NoneOrPresent              = "*"
	OpeningTag                 = "<"
	ClosingTag                 = ">"
)

type BasicParser interface {
	Parse([][2]string) error
}

type NfaGrammar struct {
	Nfa            nfa.NFA
	Target         string
	AlreadyOnFinal bool
}

type RegexpParser struct {
	nfaGrammar []*NfaGrammar
}

type RegexpParserConfig struct {
	Grammar [][2]string
}

type HandleSymbolInput struct {
	PreviousState []*nfa.State
	NfaData       *nfa.NFA
	Tag           string
	IsFinal       bool
}

func handleOneOrMore(input *HandleSymbolInput) (*nfa.NFA, []*nfa.State, error) {
	var newStates []*nfa.State
	var state1 *nfa.State
	var err error

	if input.NfaData == nil {
		input.NfaData, state1, err = nfa.NewNFA(input.Tag, false)

		if err != nil {
			return nil, nil, err
		}

	} else {

		state1, err = input.NfaData.AddState(&nfa.State{
			Name: input.Tag,
		}, false)

		if err != nil {
			return nil, nil, err
		}
	}

	state2, err := input.NfaData.AddState(&nfa.State{
		Name: input.Tag,
	}, input.IsFinal)

	if err != nil {
		return nil, nil, err
	}

	if state1 == nil || state2 == nil {
		return nil, nil, errors.New(NilState)
	}

	newStates = append(newStates, state1)
	newStates = append(newStates, state2)

	err = input.NfaData.AddTransition(state1.Index, input.Tag, *state2)
	err = input.NfaData.AddTransition(state2.Index, input.Tag, *state2)

	if err != nil {
		return nil, nil, err
	}

	return input.NfaData, newStates, nil
}

func convertGrammarToNfa(grammar string) (*nfa.NFA, error) {
	var newNFA *nfa.NFA
	var err error

	if string(grammar[0]) == "{" && string(grammar[len(grammar)-1]) == "}" {
		grammar = strings.Replace(grammar, "{", "", 1)
		grammar = strings.Replace(grammar, "}", "", 1)

		var nextTag *string
		isFinal := false
		for {
			tag := helper.GetStringInBetween(grammar, OpeningTag, ClosingTag)

			if tag == nil {
				break
			}

			processedTag := OpeningTag + *tag + ClosingTag

			grammar = strings.Replace(grammar, processedTag, "", 1)

			nextTag = helper.GetStringInBetween(grammar, "<", ">")

			if nextTag == nil {
				isFinal = true
			}

			processedCharacter := ""
			for idx, _ := range grammar {
				if string(grammar[idx]) == OneOrMore {
					newNFA, _, err = handleOneOrMore(&HandleSymbolInput{
						NfaData:       newNFA,
						Tag:           *tag,
						IsFinal:       isFinal,
						PreviousState: nil,
					})

					if err != nil {
						return nil, err
					}

					processedCharacter += OneOrMore
				} else if string(grammar[idx]) == OpeningTag {
					break
				}
			}

			grammar = strings.Replace(grammar, processedCharacter, "", 1)
		}
	} else {
		return nil, errors.New(InvalidGrammar)
	}
	return newNFA, nil
}

func NewRegexpParser(config RegexpParserConfig) (*RegexpParser, error) {
	var nfas []*NfaGrammar

	for _, g := range config.Grammar {
		newNfa, err := convertGrammarToNfa(g[1])

		if err != nil {
			return nil, err
		}

		if newNfa == nil {
			return nil, errors.New(CannotCreateNFAFromGrammar)
		}

		nfas = append(nfas, &NfaGrammar{
			Nfa:    *newNfa,
			Target: g[0],
		})
	}
	return &RegexpParser{
		nfaGrammar: nfas,
	}, nil
}

type ParsedGrammar struct {
	GeneralTag *string
	Words      [][2]string
}

func (rp *RegexpParser) ResetAllNfa() error {
	for _, nfa := range rp.nfaGrammar {
		err := nfa.Nfa.Reset()

		nfa.AlreadyOnFinal = false

		if err != nil {
			return err
		}
	}
	return nil
}

func (rp *RegexpParser) Parse(input [][2]string) ([]ParsedGrammar, error) {
	var parsedWords [][2]string

	var processedTag []string
	var processedGrammars []ParsedGrammar

	for _, word := range input {

		processedTag = append(processedTag, word[1])

		validOnPriority := false
		for _, nfa := range rp.nfaGrammar {

			res := nfa.Nfa.VerifyInputs(processedTag)
			currState, err := nfa.Nfa.GetCurrenteState()

			if err != nil {
				return nil, err
			}

			lenState := uint64(len(currState))

			if nfa.AlreadyOnFinal && lenState == 0 && !validOnPriority {
				err = rp.ResetAllNfa()

				if err != nil {
					return nil, err
				}

				tag := nfa.Target

				processedGrammars = append(processedGrammars, ParsedGrammar{
					GeneralTag: &tag,
					Words:      parsedWords,
				})

				parsedWords = nil
				processedTag = nil

				parsedWords = append(parsedWords, word)

				break
			}

			if lenState > 0 {
				validOnPriority = true
			}

			parsedWords = append(parsedWords, word)
			nfa.AlreadyOnFinal = res
		}

		if !validOnPriority {
			err := rp.ResetAllNfa()

			if err != nil {
				return nil, err
			}

			processedGrammars = append(processedGrammars, ParsedGrammar{
				Words: parsedWords,
			})

			parsedWords = nil
			processedTag = nil
		}
	}
	return processedGrammars, nil
}

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
	NoneOrMore                 = "*"
	Optional                   = "?"
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

func handleNoneOrMore(input *HandleSymbolInput) (*nfa.NFA, []*nfa.State, error) {
	var newStates []*nfa.State
	var state1 *nfa.State
	var err error

	if input.NfaData == nil {
		input.NfaData, state1, err = nfa.NewNFA(input.Tag, input.IsFinal)

		if err != nil {
			return nil, nil, err
		}

	} else {

		state1, err = input.NfaData.AddState(&nfa.State{
			Name: input.Tag,
		}, input.IsFinal)

		if err != nil {
			return nil, nil, err
		}
	}

	newStates = append(newStates, state1)

	err = input.NfaData.AddTransition(state1.Index, input.Tag, *state1)

	if err != nil {
		return nil, nil, err
	}

	for idx, _ := range input.PreviousState {
		newStates = append(newStates, input.PreviousState[idx])
		err = input.NfaData.AddTransition(input.PreviousState[idx].Index, input.Tag, *state1)

		if err != nil {
			return nil, nil, err
		}

	}

	return input.NfaData, newStates, nil
}

func handleBasic(input *HandleSymbolInput) (*nfa.NFA, []*nfa.State, error) {
	var newStates []*nfa.State
	var state1 *nfa.State
	var err error

	if input.NfaData == nil {
		input.NfaData, state1, err = nfa.NewNFA(input.Tag, input.IsFinal)

		if err != nil {
			return nil, nil, err
		}

	} else {

		state1, err = input.NfaData.AddState(&nfa.State{
			Name: input.Tag,
		}, input.IsFinal)

		if err != nil {
			return nil, nil, err
		}
	}

	newStates = append(newStates, state1)

	if err != nil {
		return nil, nil, err
	}

	for idx, _ := range input.PreviousState {
		newStates = append(newStates, input.PreviousState[idx])
		err = input.NfaData.AddTransition(input.PreviousState[idx].Index, input.Tag, *state1)

		if err != nil {
			return nil, nil, err
		}

	}

	return input.NfaData, newStates, nil
}

func handleOptional(input *HandleSymbolInput) (*nfa.NFA, []*nfa.State, error) {
	var newStates []*nfa.State
	var state1 *nfa.State
	var err error

	if input.NfaData == nil {
		input.NfaData, state1, err = nfa.NewNFA(input.Tag, input.IsFinal)

		if err != nil {
			return nil, nil, err
		}

	} else {

		state1, err = input.NfaData.AddState(&nfa.State{
			Name: input.Tag,
		}, input.IsFinal)

		if err != nil {
			return nil, nil, err
		}
	}

	newStates = append(newStates, state1)

	err = input.NfaData.AddTransition(state1.Index, input.Tag, *state1)

	if err != nil {
		return nil, nil, err
	}

	for idx, _ := range input.PreviousState {
		newStates = append(newStates, input.PreviousState[idx])
		err = input.NfaData.AddTransition(input.PreviousState[idx].Index, input.Tag, *state1)

		if err != nil {
			return nil, nil, err
		}

	}

	return input.NfaData, newStates, nil
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

	for idx, _ := range input.PreviousState {
		err = input.NfaData.AddTransition(input.PreviousState[idx].Index, input.Tag, *state1)

		if err != nil {
			return nil, nil, err
		}

	}

	newStates = append(newStates, state2)

	err = input.NfaData.AddTransition(state1.Index, input.Tag, *state2)

	if err != nil {
		return nil, nil, err
	}

	err = input.NfaData.AddTransition(state2.Index, input.Tag, *state2)

	if err != nil {
		return nil, nil, err
	}

	return input.NfaData, newStates, nil
}

func convertGrammarToNfa(grammar string) (*nfa.NFA, error) {
	var newNFA *nfa.NFA
	var err error
	var previousState []*nfa.State

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

			if len(grammar) == 0 {
				newNFA, previousState, err = handleBasic(&HandleSymbolInput{
					NfaData:       newNFA,
					Tag:           *tag,
					IsFinal:       isFinal,
					PreviousState: previousState,
				})

				if err != nil {
					return nil, err
				}

				break
			}
			for idx, _ := range grammar {
				if string(grammar[idx]) == OneOrMore {
					newNFA, previousState, err = handleOneOrMore(&HandleSymbolInput{
						NfaData:       newNFA,
						Tag:           *tag,
						IsFinal:       isFinal,
						PreviousState: previousState,
					})

					if err != nil {
						return nil, err
					}

					processedCharacter += OneOrMore
				} else if string(grammar[idx]) == Optional {
					newNFA, previousState, err = handleOptional(&HandleSymbolInput{
						NfaData:       newNFA,
						Tag:           *tag,
						IsFinal:       isFinal,
						PreviousState: previousState,
					})

					if err != nil {
						return nil, err
					}

					processedCharacter += Optional
				} else if string(grammar[idx]) == NoneOrMore {
					newNFA, previousState, err = handleNoneOrMore(&HandleSymbolInput{
						NfaData:       newNFA,
						Tag:           *tag,
						IsFinal:       isFinal,
						PreviousState: previousState,
					})

					if err != nil {
						return nil, err
					}

					processedCharacter += NoneOrMore
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

		newNfa.PrintTransitionTable()
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
	for _, nfaGrammar := range rp.nfaGrammar {
		err := nfaGrammar.Nfa.Reset()

		nfaGrammar.AlreadyOnFinal = false

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

	for idxWord, word := range input {

		processedTag = append(processedTag, word[1])

		validOnPriority := false
		for _, nfaGrammar := range rp.nfaGrammar {

			res := nfaGrammar.Nfa.VerifyInputs(processedTag)

			processedTag = nil
			currState, err := nfaGrammar.Nfa.GetCurrenteState()

			if err != nil {
				return nil, err
			}

			lenState := uint64(len(currState))

			if nfaGrammar.AlreadyOnFinal && lenState == 0 && !validOnPriority {
				err = rp.ResetAllNfa()

				if err != nil {
					return nil, err
				}

				tag := nfaGrammar.Target

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
			nfaGrammar.AlreadyOnFinal = res
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

		if idxWord == len(input)-1 {
			found := false
			for _, nfaGrammar := range rp.nfaGrammar {
				tag := nfaGrammar.Target

				processedGrammars = append(processedGrammars, ParsedGrammar{
					GeneralTag: &tag,
					Words:      parsedWords,
				})

				found = true
				break
			}

			if !found {
				processedGrammars = append(processedGrammars, ParsedGrammar{
					Words: parsedWords,
				})
			}
		}
	}
	return processedGrammars, nil
}

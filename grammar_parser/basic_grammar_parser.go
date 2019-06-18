package grammar_parser

import (
	"errors"
	"github.com/adrian3ka/go-learn-ai/helper"
	"github.com/adrian3ka/go-learn-ai/nfa"
	"strings"
)

//AVAIABLE SYMBOL
//- CHINCKING : } {
//- CHUNGKING : { }

const (
	CannotCreateNFAFromGrammar    = "Cannot Create NFA From Grammar"
	InvalidGrammar                = "Invalid Grammar"
	NilState                      = "Nil State"
	ChinkingInitialStateMustBeNil = "Chinking Initial State Must Be Nil"
	OneOrMore                     = "+"
	NoneOrMore                    = "*"
	Optional                      = "?"
	OpeningTag                    = "<"
	ClosingTag                    = ">"

	Chunking = "Chunking"
	Chinking = "Chinking"
)

type BasicParser interface {
	Parse([][2]string) error
}

type NfaGrammar struct {
	Nfa            nfa.NFA
	Target         string
	AlreadyOnFinal bool
	Type           string
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

	var tags []string
	tags = strings.Split(input.Tag, "|")

	for _, tag := range tags {
		err = input.NfaData.AddTransition(state1.Index, tag, *state1)
	}

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

func handleInitialChunking(input *HandleSymbolInput) (*nfa.NFA, []*nfa.State, error) {
	var newStates []*nfa.State
	var state1 *nfa.State
	var err error

	if input.NfaData == nil {
		input.NfaData, state1, err = nfa.NewNFA(nfa.Negate+input.Tag, false)

		if err != nil {
			return nil, nil, err
		}

	} else {
		return nil, nil, errors.New(ChinkingInitialStateMustBeNil)
	}

	var tags []string
	tags = strings.Split(input.Tag, "|")

	for _, tag := range tags {
		err = input.NfaData.AddTransition(state1.Index, nfa.Negate+tag, *state1)
	}

	newStates = append(newStates, state1)

	return input.NfaData, newStates, nil
}

func handleEndOfChunking(input *HandleSymbolInput) (*nfa.NFA, []*nfa.State, error) {
	var newStates []*nfa.State
	var state1 *nfa.State
	var err error

	if input.NfaData == nil {
		input.NfaData, state1, err = nfa.NewNFA(nfa.Negate+input.Tag, true)

		if err != nil {
			return nil, nil, err
		}

	} else {
		state1, err = input.NfaData.AddState(&nfa.State{
			Name: nfa.Negate + input.Tag,
		}, true)

		if err != nil {
			return nil, nil, err
		}
	}

	var tags []string
	tags = strings.Split(input.Tag, "|")

	for _, tag := range tags {
		err = input.NfaData.AddTransition(state1.Index, nfa.Negate+tag, *state1)
	}

	newStates = append(newStates, state1)

	for idx, _ := range input.PreviousState {
		err = input.NfaData.AddTransition(input.PreviousState[idx].Index, input.Tag, *state1)

		if err != nil {
			return nil, nil, err
		}
	}

	return input.NfaData, newStates, nil
}

func convertGrammarToNfa(grammar string) (*nfa.NFA, *string, error) {
	var newNFA *nfa.NFA
	var grammarType string
	var err error
	var previousState []*nfa.State
	var processedTag string
	var tag *string

	var nextTag *string
	isFinal := false

	if string(grammar[0]) == "{" && string(grammar[len(grammar)-1]) == "}" {
		grammarType = Chunking
	} else if string(grammar[0]) == "}" && string(grammar[len(grammar)-1]) == "{" {
		grammarType = Chinking
	} else {
		return nil, nil, errors.New(InvalidGrammar)
	}

	if grammarType == Chunking {
		grammar = strings.Replace(grammar, "{", "", 1)
		grammar = strings.Replace(grammar, "}", "", 1)
	} else if grammarType == Chinking {
		grammar = strings.Replace(grammar, "}", "", 1)
		grammar = strings.Replace(grammar, "{", "", 1)

		tag := helper.GetStringInBetween(grammar, OpeningTag, ClosingTag)

		if tag == nil {
			return nil, nil, nil
		}

		nextTag = helper.GetStringInBetween(grammar, "<", ">")
		isFinal = false

		newNFA, previousState, err = handleInitialChunking(&HandleSymbolInput{
			NfaData:       newNFA,
			Tag:           *tag,
			IsFinal:       isFinal,
			PreviousState: previousState,
		})

		if err != nil {
			return nil, nil, err
		}
	}

	for {
		tag = helper.GetStringInBetween(grammar, OpeningTag, ClosingTag)

		if tag == nil {
			break
		}

		processedTag = OpeningTag + *tag + ClosingTag

		grammar = strings.Replace(grammar, processedTag, "", 1)

		nextTag = helper.GetStringInBetween(grammar, "<", ">")

		if nextTag == nil && grammarType != Chinking {
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
				return nil, nil, err
			}

			newNFA.PrintTransitionTable()
			break
		}
		processed := false
		for idx, _ := range grammar {
			if string(grammar[idx]) == OneOrMore {
				newNFA, previousState, err = handleOneOrMore(&HandleSymbolInput{
					NfaData:       newNFA,
					Tag:           *tag,
					IsFinal:       isFinal,
					PreviousState: previousState,
				})

				if err != nil {
					return nil, nil, err
				}

				processedCharacter += OneOrMore
				processed = true
			} else if string(grammar[idx]) == Optional {
				newNFA, previousState, err = handleOptional(&HandleSymbolInput{
					NfaData:       newNFA,
					Tag:           *tag,
					IsFinal:       isFinal,
					PreviousState: previousState,
				})

				if err != nil {
					return nil, nil, err
				}

				processedCharacter += Optional
				processed = true
			} else if string(grammar[idx]) == NoneOrMore {
				newNFA, previousState, err = handleNoneOrMore(&HandleSymbolInput{
					NfaData:       newNFA,
					Tag:           *tag,
					IsFinal:       isFinal,
					PreviousState: previousState,
				})

				if err != nil {
					return nil, nil, err
				}

				processedCharacter += NoneOrMore
				processed = true
			} else if string(grammar[idx]) == OpeningTag {
				if !processed {
					newNFA, previousState, err = handleBasic(&HandleSymbolInput{
						NfaData:       newNFA,
						Tag:           *tag,
						IsFinal:       isFinal,
						PreviousState: previousState,
					})

					if err != nil {
						return nil, nil, err
					}
				}
				break
			}
		}

		grammar = strings.Replace(grammar, processedCharacter, "", 1)
	}

	if grammarType == Chinking {
		newNFA, previousState, err = handleEndOfChunking(&HandleSymbolInput{
			NfaData:       newNFA,
			Tag:           *tag,
			IsFinal:       true,
			PreviousState: previousState,
		})

		if err != nil {
			return nil, nil, err
		}
	}

	return newNFA, &grammarType, nil
}

func NewRegexpParser(config RegexpParserConfig) (*RegexpParser, error) {
	var nfas []*NfaGrammar

	for _, g := range config.Grammar {
		newNfa, grammarType, err := convertGrammarToNfa(g[1])

		if err != nil {
			return nil, err
		}

		if newNfa == nil {
			return nil, errors.New(CannotCreateNFAFromGrammar)
		}

		newNfa.PrintTransitionTable()
		nfas = append(nfas, &NfaGrammar{
			Nfa:    *newNfa,
			Type:   *grammarType,
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
	for idx, _ := range rp.nfaGrammar {
		err := rp.nfaGrammar[idx].Nfa.Reset()

		rp.nfaGrammar[idx].AlreadyOnFinal = false

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
		for idx, _ := range rp.nfaGrammar {

			res := rp.nfaGrammar[idx].Nfa.VerifyInputs(processedTag)

			currState, err := rp.nfaGrammar[idx].Nfa.GetCurrenteState()

			if err != nil {
				return nil, err
			}

			lenState := uint64(len(currState))

			if rp.nfaGrammar[idx].AlreadyOnFinal && lenState == 0 && !validOnPriority {
				err = rp.ResetAllNfa()

				if err != nil {
					return nil, err
				}

				tag := rp.nfaGrammar[idx].Target

				processedGrammars = append(processedGrammars, ParsedGrammar{
					GeneralTag: &tag,
					Words:      parsedWords,
				})

				parsedWords = nil
				processedTag = nil

				break
			}

			if lenState > 0 && !validOnPriority {
				validOnPriority = true
			}

			rp.nfaGrammar[idx].AlreadyOnFinal = res
		}

		processedTag = nil
		parsedWords = append(parsedWords, word)

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

		if idxWord == len(input)-1 && len(parsedWords) > 0 {
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

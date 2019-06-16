package nfa

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	StateNotFound = "State Not Found"
	InvalidInput  = "Invalid Input"

	Negate = "!"
)

type transitionInput struct {
	srcStateIndex uint64
	input         string
}

type State struct {
	Name  string
	Index uint64
}

type destState map[State]bool

type NFA struct {
	initState    State
	currentState map[State]bool
	allStates    []State
	finalStates  []State
	transition   map[transitionInput]destState
	inputMap     map[string]bool
}

//New a new NFA
func NewNFA(initStateName string, isFinal bool) (*NFA, *State, error) {
	initState := State{
		Name:  initStateName,
		Index: 0,
	}
	retNFA := &NFA{
		transition: make(map[transitionInput]destState),
		inputMap:   make(map[string]bool),
		initState:  initState,
	}

	retNFA.currentState = make(map[State]bool)
	retNFA.currentState[initState] = true
	_, err := retNFA.AddState(&initState, isFinal)

	if err != nil {
		return nil, nil, err
	}

	return retNFA, &initState, nil
}

func (d *NFA) GetCurrenteState() (map[State]bool, error) {
	return d.currentState, nil
}
func (d *NFA) GetAllState() ([]State, error) {
	return d.allStates, nil
}

//Add new state in this NFA
func (d *NFA) AddState(state *State, isFinal bool) (*State, error) {
	if state == nil || state.Name == "" {
		return nil, errors.New(InvalidInput)
	}

	currentIndex := uint64(len(d.allStates))

	state.Index = currentIndex

	d.allStates = append(d.allStates, *state)
	if isFinal {
		d.finalStates = append(d.finalStates, *state)
	}

	return state, nil
}

//Add new transition function into NFA
func (d *NFA) AddTransition(srcStateIndex uint64, input string, dstStateList ...State) error {
	find := false

	for _, v := range d.allStates {
		if v.Index == srcStateIndex {
			find = true
		}
	}

	if !find {
		return errors.New(StateNotFound)
	}

	if input == "" {
		return errors.New(InvalidInput)
	}

	//find input if exist in NFA input List
	if _, ok := d.inputMap[input]; !ok {
		//not exist, new input in this NFA
		d.inputMap[input] = true
	}

	dstMap := make(map[State]bool)
	for _, destState := range dstStateList {
		dstMap[destState] = true
	}

	targetTrans := transitionInput{srcStateIndex: srcStateIndex, input: input}
	d.transition[targetTrans] = dstMap

	return nil
}

func (d *NFA) Input(testInput string) ([]State, error) {
	updateCurrentState := make(map[State]bool)
	for current, _ := range d.currentState {
		intputTrans := transitionInput{srcStateIndex: current.Index, input: testInput}

		if valMap, ok := d.transition[intputTrans]; ok {
			for dst, _ := range valMap {
				updateCurrentState[dst] = true
			}
		} else {
			//dead state, remove in current state
			//do nothing.
			for key, _ := range d.transition {
				var temp string
				if string(key.input[0]) != Negate {
					continue
				}
				temp = strings.Replace(key.input, Negate, "", 1)
				if temp != testInput {
					for dst, _ := range d.transition[key] {
						updateCurrentState[dst] = true
					}
				}
			}
		}
	}

	//update curret state
	d.currentState = updateCurrentState

	//return result
	var ret []State
	for state, _ := range updateCurrentState {
		ret = append(ret, state)
	}

	return ret, nil
}

//To verify current state if it is final state
func (d *NFA) Verify() bool {
	for _, val := range d.finalStates {
		for cState, _ := range d.currentState {
			if val == cState {
				return true
			}
		}
	}
	return false
}

//Reset NFA state to initilize state, but all state and transition function will remain
func (d *NFA) Reset() error {
	initState := make(map[State]bool)
	initState[d.initState] = true
	d.currentState = initState

	return nil
}

//Verify if list of input could be accept by NFA or not
func (d *NFA) VerifyInputs(inputs []string) bool {
	for _, v := range inputs {
		d.Input(v)
	}
	return d.Verify()
}

func (d *NFA) PrintTransitionTable() {
	fmt.Println("==========================================================================================================================")
	//list all inputs
	var inputList []string

	fmt.Printf("%16s|", "")
	for key, _ := range d.inputMap {
		fmt.Printf("%15s|", key)
		inputList = append(inputList, key)
	}

	fmt.Printf("\n")
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------")

	for _, state := range d.allStates {
		isFinal := false

		for _, s := range d.finalStates {
			if s.Index == state.Index {
				isFinal = true
				break
			}
		}
		if isFinal {
			fmt.Printf("*%4d %10s|", state.Index, state.Name)
		} else {
			fmt.Printf("%5d %10s|", state.Index, state.Name)
		}
		for _, key := range inputList {
			checkInput := transitionInput{srcStateIndex: state.Index, input: key}
			if dstState, ok := d.transition[checkInput]; ok {
				var temp []string

				for val, _ := range dstState {
					temp = append(temp, strconv.FormatUint(val.Index, 10))
				}
				fmt.Printf("%15s|", strings.Join(temp, ","))
			} else {
				fmt.Printf("%15s|", "NA")
			}
		}
		fmt.Printf("\n")
	}

	fmt.Println("--------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("==========================================================================================================================")
}

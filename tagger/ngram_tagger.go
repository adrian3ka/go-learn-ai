package tagger

import (
	"container/list"
	"github.com/adrian3ka/go-learn-ai/helper"
	"strings"
)

type UnigramTagger struct {
	mapTag        map[string]string
	backoffTagger Tagger
}

type UnigramTaggerConfig struct {
	BackoffTagger Tagger
}

func NewUnigramTagger(cfg UnigramTaggerConfig) *UnigramTagger {
	u := UnigramTagger{
		backoffTagger: cfg.BackoffTagger,
	}
	u.mapTag = make(map[string]string)
	return &u
}

func (u *UnigramTagger) Predict(text string) ([][2]string, error) {
	splitedStrings := strings.Split(text, " ")
	var tuple [][2]string

	for _, splitedString := range splitedStrings {
		var selectedTag *string

		if helper.IsLetter(splitedString) {
			if val, exists := u.mapTag[splitedString]; exists {
				selectedTag = &val
			}
		}

		if selectedTag == nil && u.backoffTagger != nil {
			predictedValue, err := u.backoffTagger.Predict(splitedString)

			if err != nil {
				return nil, err
			}

			selectedTag = &predictedValue[0][1]
		}

		if selectedTag == nil {
			x := ""
			selectedTag = &x
		}

		tuple = append(tuple, [2]string{
			splitedString,
			*selectedTag,
		})
	}
	return tuple, nil
}

func (u *UnigramTagger) Learn(tuple [][][2]string) error {
	var tupleMap = make(map[string]map[string]float64)

	for _, sentence := range tuple {
		for _, word := range sentence {
			if _, exists := tupleMap[word[0]]; !exists {
				var temp = make(map[string]float64)
				temp[word[1]] = 1
				tupleMap[word[0]] = temp
			} else {
				if _, exists := tupleMap[word[0]][word[1]]; !exists {
					tupleMap[word[0]][word[1]] = 1
				} else {
					tupleMap[word[0]][word[1]] += 1
				}
			}
		}
	}

	for word, tm := range tupleMap {
		selectedTag := ""
		maxCount := float64(0)
		for tag, count := range tm {
			if maxCount < count {
				selectedTag = tag
				maxCount = count
			}
		}
		u.mapTag[word] = selectedTag
	}

	return nil
}

type NGramTagger struct {
	mapTag        map[string]string
	n             uint64
	backoffTagger Tagger
}

type NGramTaggerConfig struct {
	BackoffTagger Tagger
	N             uint64
}

func NewNGramTagger(cfg NGramTaggerConfig) *NGramTagger {
	if cfg.N < 2 {
		cfg.N = 2
	}
	n := NGramTagger{
		backoffTagger: cfg.BackoffTagger,
		n:             cfg.N,
	}
	n.mapTag = make(map[string]string)
	return &n
}

func (n *NGramTagger) Predict(text string) ([][2]string, error) {
	splitedStrings := strings.Split(text, " ")
	var tuple [][2]string
	minimumWord := n.n - 1

	for idx, splitedString := range splitedStrings {
		var selectedTag *string

		if uint64(idx) >= minimumWord && helper.IsLetter(splitedString) {

			generatedTag := ""

			for i := 1; i <= int(minimumWord); i++ {
				generatedTag += tuple[len(tuple)-i][1] + "_"
			}

			generatedTag += splitedString

			if helper.IsAlphaUnderscore(generatedTag) {
				if val, exists := n.mapTag[generatedTag]; exists {
					if helper.IsLetter(val) {
						selectedTag = &val
					}
				}
			}
		}

		if selectedTag == nil && n.backoffTagger != nil {
			predictedValue, err := n.backoffTagger.Predict(splitedString)

			if err != nil {
				return nil, err
			}

			selectedTag = &predictedValue[0][1]
		}

		if selectedTag == nil {
			x := ""
			selectedTag = &x
		}

		tuple = append(tuple, [2]string{
			splitedString,
			*selectedTag,
		})
	}
	return tuple, nil
}

func (n *NGramTagger) Learn(tuple [][][2]string) error {
	var tupleMap = make(map[string]map[string]float64)
	maxQueueCount := n.n - 1
	generatedTag := ""
	for _, sentence := range tuple {

		queue := list.New()
		for idx, word := range sentence {
			if idx == len(sentence)-1 {
				break
			}

			if uint64(idx) < n.n-1 {
				queue.PushBack(word)
				continue
			}

			generatedTag = ""

			for i := 0; i < queue.Len(); i++ {
				e := queue.Front()
				val := e.Value.([2]string)
				generatedTag += val[1] + "_"
				queue.MoveToBack(e)
			}

			generatedTag += word[0]

			if _, exists := tupleMap[generatedTag]; !exists {
				var temp = make(map[string]float64)
				temp[word[1]] = 1
				tupleMap[generatedTag] = temp
			} else {
				if _, exists := tupleMap[generatedTag][word[1]]; !exists {
					tupleMap[generatedTag][word[1]] = 1
				} else {
					tupleMap[generatedTag][word[1]] += 1
				}
			}

			queue.PushBack(word)

			if uint64(queue.Len()) > maxQueueCount {
				e := queue.Front()
				queue.Remove(e)
			}
		}
	}

	for word, tm := range tupleMap {
		selectedTag := ""
		maxCount := float64(0)
		for tag, count := range tm {
			if maxCount < count {
				selectedTag = tag
				maxCount = count
			}
		}
		n.mapTag[word] = selectedTag
	}

	return nil
}

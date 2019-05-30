package tagger

import "strings"

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

		if val, exists := u.mapTag[splitedString]; exists {
			selectedTag = &val
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

func (u *UnigramTagger) Learn(tuple [][2]string) error {
	var tupleMap = make(map[string]map[string]float64)

	for _, t := range tuple {
		if _, exists := tupleMap[t[0]]; !exists {
			var temp = make(map[string]float64)
			temp[t[1]] = 1
			tupleMap[t[0]] = temp
		} else {
			if _, exists := tupleMap[t[0]][t[1]]; !exists {
				tupleMap[t[0]][t[1]] = 1
			} else {
				tupleMap[t[0]][t[1]] += 1
			}
		}
	}

	for word, tm := range tupleMap {
		selectedTag := ""
		maxCount := float64(0)
		for tag, count := range tm {
			if maxCount < count {
				selectedTag = tag
			}
		}
		u.mapTag[word] = selectedTag
	}

	return nil
}

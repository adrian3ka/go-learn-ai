package word_embedding

type GetNeighborWordsInput struct {
	Corpus     []string
	WindowSize uint64
}

type NeighborhoodWord struct {
	Words    [2]string
	Distance uint64
}
type GetNeighborWordsOutput struct {
	NeighborhoodWords []NeighborhoodWord
}

func GetNeighborWords(input GetNeighborWordsInput) GetNeighborWordsOutput {
	var output GetNeighborWordsOutput

	return output
}

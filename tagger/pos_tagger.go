package tagger

import (
	"github.com/adrian3ka/go-learn-ai/helper"
	"regexp"
	"strings"
)

var Punctuation = []string{
	".", "?", "!", ",", "-", "--",
}

var IndonesianStopWords = []string{
	"dan", "di", "ke", "dari", "ok", "ya", "pula", "pada", "ke", "yang", "ia",
}

var DefaultSimpleIndonesianRegexTagger = [][2]string{
	[2]string{`^[0-9]*$`, `CDP`},
	[2]string{`^(satu|dua|tiga|empat|lima|enam|tujuh|delapan|sembilan|sepuluh)*$`, `CDP`},
	[2]string{`(bidang)$`, `VB`},
	[2]string{`^di.+(kan|i)$`, `VB`},
	[2]string{`(membuat)$`, `VBT`},
	[2]string{`(tidak|tak)$`, `NEG`},
	[2]string{`se(baik|benar|tidak|layak|lekas|sungguh|yogya|belum|pantas|balik|lanjut)(nya)$`, `RB`},
	[2]string{`(sekadar|amat|bahkan|cukup|jua|justru|kembali|kurang|malah|mau|nian|niscaya|pasti|patut|perlu|lagi|pernah|pun|sekali|selalu|senantiasa|sering|sungguh|tentu|terus|lebih|hampir|jarang|juga|kerap|makin|memang|nyaris|paling|pula|saja|saling|sangat|segera|semakin|serba|entah|hanya|kadangkala)$`, `RB`},
	[2]string{`(akan|antara|bagi|buat|dari|dengan|di|ke|kecuali|lepas|oleh|pada|per|peri|seperti|tanpa|tentang|untuk|dengan)$`, `IN`},
	[2]string{`(dan|serta|atau|tetapi|melainkan|padahal|sedangkan)$`, `CC`},
	[2]string{`(yang|sejak|semenjak|sedari|sewaktu|ketika|tatkala|sementara|begitu|seraya|selagi|selama|serta|sambil|demi|setelah|sesudah|sebelum|sehabis|selesai|seusai|hingga|sampai|jika|kalau|jikalau|asal)$`, `SC`},
}

type StringToTupleInput struct {
	Text        string
	Lower       bool
	Punctuation []string
	Simplify    bool
	Default     *string
}

type StringToTupleOutput struct {
	Tuple       [][][2]string
	Punctuation []string
}

func TupleSplit(r rune) bool {
	return r == '\n' || r == ' '
}

func StringToTuple(input StringToTupleInput) StringToTupleOutput {
	if len(input.Punctuation) == 0 {
		input.Punctuation = Punctuation
	}

	splitedStrings := strings.FieldsFunc(input.Text, TupleSplit)
	var tuple [][][2]string

	var tempSentence [][2]string
	for _, splitedString := range splitedStrings {
		var splittedWordAndTag [2]string
		temp := helper.LastSplit(splitedString, "/")

		if len(temp) != 2 {
			continue
		}

		if input.Lower {
			splittedWordAndTag[0] = strings.ToLower(temp[0])
		} else {
			splittedWordAndTag[0] = temp[0]
		}

		tag := strings.ToUpper(temp[1])

		if len(tag) > 3 && input.Default != nil {
			tag = *input.Default
		}

		splittedWordAndTag[1] = tag

		tempSentence = append(tempSentence, splittedWordAndTag)

		if helper.IsStringEqual(splittedWordAndTag[0], input.Punctuation) {
			tuple = append(tuple, tempSentence)
			tempSentence = nil
		}

	}

	return StringToTupleOutput{
		Tuple: tuple,
	}
}

type Tagger interface {
	Learn(tuple [][][2]string) error
	Predict(text string) ([][2]string, error)
}

type DefaultTagger struct {
	defaultTag string
}

type DefaultTaggerConfig struct {
	DefaultTag string
}

func NewDefaultTagger(cfg DefaultTaggerConfig) *DefaultTagger {
	return &DefaultTagger{
		defaultTag: strings.ToUpper(cfg.DefaultTag),
	}
}

func (n *DefaultTagger) Learn(tuple [][][2]string) error {
	return nil
}

func (n *DefaultTagger) Predict(text string) ([][2]string, error) {
	splitedStrings := strings.Split(text, " ")
	var tuple [][2]string

	for _, splitedString := range splitedStrings {
		tag := n.defaultTag

		if !helper.IsAlphaNumeric(splitedString) && len(splitedString) == 1 {
			tag = splitedString
		}
		tuple = append(tuple, [2]string{
			splitedString,
			tag,
		})
	}
	return tuple, nil
}

type CompliedPattern struct {
	Pattern *regexp.Regexp
	Tag     string
}

type RegexTagger struct {
	compliedPatterns []CompliedPattern
	backoffTagger    Tagger
}

type RegexTaggerConfig struct {
	Patterns      [][2]string
	BackoffTagger Tagger
}

func NewRegexTagger(cfg RegexTaggerConfig) *RegexTagger {
	var compliedPatterns []CompliedPattern
	for _, pattern := range cfg.Patterns {
		cp := CompliedPattern{
			Pattern: regexp.MustCompile(pattern[0]),
			Tag:     pattern[1],
		}
		compliedPatterns = append(compliedPatterns, cp)
	}
	return &RegexTagger{
		compliedPatterns: compliedPatterns,
		backoffTagger:    cfg.BackoffTagger,
	}
}

func (n *RegexTagger) Learn(tuple [][][2]string) error {
	return nil
}

func (n *RegexTagger) Predict(text string) ([][2]string, error) {
	splitedStrings := strings.Split(text, " ")
	var tuple [][2]string
	for _, splitedString := range splitedStrings {
		var tag *string
		for _, compiledPattern := range n.compliedPatterns {
			if compiledPattern.Pattern.MatchString(splitedString) {
				x := compiledPattern.Tag
				tag = &x
				break
			}
		}
		if n.backoffTagger != nil && tag == nil {
			result, err := n.backoffTagger.Predict(splitedString)

			if err != nil {
				return nil, err
			}

			tag = &result[0][1]
		}

		if tag == nil {
			x := ""
			tag = &x
		}
		tuple = append(tuple, [2]string{
			splitedString,
			*tag,
		})
	}
	return tuple, nil
}

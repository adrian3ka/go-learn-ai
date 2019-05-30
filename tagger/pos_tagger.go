package tagger

import (
	"regexp"
	"strings"
)

var DefaultIndonesianRegexTagger = [][2]string{
	[2]string{`^-?[0-9]+(.[0-9]+)?$`, `cdp`},
	[2]string{`meng[aiueokghx].+$`, `vb`},
	[2]string{`mem[bpf]([a-df-z][a-qs-z]|er).+$`, `vb`},
	[2]string{`me[lnryw](a-df-z).+$`, `vb`},
	[2]string{`men[dtcjys].+$`, `vb`},
	[2]string{`di.+(kan|i)$`, `vb`},
	[2]string{`per.+(kan|i|.)$`, `vb`},
	[2]string{`ber.+(kan|an|.)$`, `vb`},
	[2]string{`(bidang)$`, `vb`},
	[2]string{`ter.+(kan|i|.)$`, `vb`},
	[2]string{`ke.+(i|an)$`, `vb`},
	[2]string{`se(baik|benar|tidak|layak|lekas|sungguh|yogya|belum|pantas|balik|lanjut)(nya)$`, `rb`},
	[2]string{`(sekadar|amat|bahkan|cukup|jua|justru|kembali|kurang|malah|mau|nian|niscaya|pasti|patut|perlu|lagi|pernah|pun|sekali|selalu|senantiasa|sering|sungguh|tentu|terus|lebih|hampir|jarang|juga|kerap|makin|memang|nyaris|paling|pula|saja|saling|sangat|segera|semakin|serba|entah|hanya|kadangkala)$`, `rb`},
	[2]string{`(akan|antara|bagi|buat|dari|dengan|di|ke|kecuali|lepas|oleh|pada|per|peri|seperti|tanpa|tentang|untuk)$`, `in`},
	[2]string{`(dan|serta|atau|tetapi|melainkan|padahal|sedangkan)$`, `cc`},
	[2]string{`(sejak|semenjak|sedari|sewaktu|ketika|tatkala|sementara|begitu|seraya|selagi|selama|serta|sambil|demi|setelah|sesudah|sebelum|sehabis|selesai|seusai|hingga|sampai|jika|kalau|jikalau|asal)$`, `sc`},
}

type StringToTupleInput struct {
	Text  string
	Lower bool
}

type StringToTupleOutput struct {
	Tuple [][2]string
}

func StringToTuple(input StringToTupleInput) StringToTupleOutput {
	splitedStrings := strings.Split(input.Text, " ")
	var tuple [][2]string

	regexReplacers := [][2]string{
		{`\s+`, ``},
	}

	for _, splitedString := range splitedStrings {
		var splittedWordAndTag [2]string
		temp := strings.Split(splitedString, "/")

		if len(temp) != 2 {
			continue
		}

		for _, regexReplacer := range regexReplacers {
			reg, err := regexp.Compile(regexReplacer[0])

			if err != nil {
				continue
			}

			temp[1] = reg.ReplaceAllString(temp[1], regexReplacer[1])
			temp[0] = reg.ReplaceAllString(temp[0], regexReplacer[1])
		}

		if input.Lower {
			splittedWordAndTag[0] = strings.ToLower(temp[0])
		}

		splittedWordAndTag[1] = strings.ToLower(temp[1])

		tuple = append(tuple, splittedWordAndTag)
	}

	return StringToTupleOutput{
		Tuple: tuple,
	}
}

type Tagger interface {
	Learn(tuple [][2]string) error
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
		defaultTag: cfg.DefaultTag,
	}
}

func (n *DefaultTagger) Learn(tuple [][2]string) error {
	return nil
}

func (n *DefaultTagger) Predict(text string) ([][2]string, error) {
	splitedStrings := strings.Split(text, " ")
	var tuple [][2]string

	for _, splitedString := range splitedStrings {
		tuple = append(tuple, [2]string{
			splitedString,
			n.defaultTag,
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

func (n *RegexTagger) Learn(tuple [][2]string) error {
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

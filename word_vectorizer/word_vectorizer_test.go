package word_vectorizer

import (
	"testing"
)

var (
	ExpectedDictionaryLength = 31
	ClassCount               = 3
	PulsaDocumentCount       = 5
	TiketDocumentCount       = 4
	SaldoDocumentCount       = 4

	Pulsa = "pulsa"
	Tiket = "tiket"
	Saldo = "saldo"
)

func TestBasic(t *testing.T) {
	wordVectorizer := New(WordVectorizerConfig{
		Lower: true,
	})

	var corpuses map[string][]string

	corpuses = make(map[string][]string)

	corpuses[Pulsa] = []string{
		"Saya mau beli pulsa dong. Jual voucher gak bang?. Mau isi pulsa dong.",
		"jual pulsa gak ya?",
		"kamu jual voucher ga?",
		"mau isi paket data bisa?",
		"mau isi pulsa bisa ga ya?",
	}

	corpuses[Tiket] = []string{
		"kamu jual tiket pesawat ga?",
		"disini jual tiket ga ya?",
		"bisa beli tiket    kereta?",
		"jual tiket apa  ya?",
	}

	corpuses[Saldo] = []string{
		"halo aku mau isi saldo dong",
		"eh mau topup dong bisa ga?",
		"mau nambah saldo dong bisa gak",
		"tolong bantu isi saldo dong 50 ribu",
	}

	err := wordVectorizer.Learn(corpuses)

	if err != nil {
		panic(err)
	}

	if ExpectedDictionaryLength != len(wordVectorizer.GetVectorizedWord()) {
		t.Errorf("Vectorized Word Length Should Be %d", ExpectedDictionaryLength)
	}

	if ClassCount != len(wordVectorizer.GetCleanedCorpus()) {
		t.Errorf("Corpus Class Should Be Still %d", ClassCount)
	}

	if len(corpuses[Pulsa]) != PulsaDocumentCount {
		t.Errorf("Pulsa Document Count Should Be Still %d", PulsaDocumentCount)
	}

	if len(corpuses[Tiket]) != TiketDocumentCount {
		t.Errorf("Tiket Document Count Should Be Still %d", TiketDocumentCount)
	}

	if len(corpuses[Saldo]) != SaldoDocumentCount {
		t.Errorf("Saldo Document Count Should Be Still %d", SaldoDocumentCount)
	}
}

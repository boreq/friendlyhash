package friendlyhash

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"
)

func ExampleNew() {
	// Create
	dictionary := []string{
		"word1",
		"word2",
		"word3",
		"word4",
		"word5",
		"word6",
	}

	h, err := New(dictionary, 2)
	if err != nil {
		panic(err)
	}

	// Humanize
	humanized, err := h.Humanize([]byte{'a', 'b'})
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.Join(humanized, "-"))

	// Dehumanize
	dehumanized, err := h.Dehumanize(humanized)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", dehumanized)

	// Output:
	// word2-word3-word1-word2-word2-word3-word1-word3
	// "ab"
}

var wordsTestCases = []struct {
	words         []string
	errorExpected bool
	message       string
}{
	{
		words:         nil,
		errorExpected: true,
		message:       "words can't be nil",
	},
	{
		words:         []string{},
		errorExpected: true,
		message:       "words can't be empty",
	},
	{
		words:         []string{"a"},
		errorExpected: true,
		message:       "words can't contain one value",
	},
	{
		words:         []string{"a", "b"},
		errorExpected: false,
		message:       "words can contain two values",
	},
	{
		words:         []string{"a", "b", "a"},
		errorExpected: true,
		message:       "words can't contain duplicates",
	},
}

func TestNewWords(t *testing.T) {
	for _, testCase := range wordsTestCases {
		if _, err := New(testCase.words, 0); err == nil {
			if testCase.errorExpected {
				t.Fatalf("%s: expected an error, got nil", testCase.message)
			}
		} else {
			if !testCase.errorExpected {
				t.Fatalf("%s: expected nil, got %s", testCase.message, err)
			}
		}
	}
}

func TestHumanizeInvalidLength(t *testing.T) {
	h, err := New([]string{"a", "b"}, 2)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	if _, err := h.Humanize([]byte{1}); err == nil {
		t.Fatal("expected an eror, got nil")
	}
}

func TestDehumanizeInvalidLength(t *testing.T) {
	h, err := New([]string{"a", "b"}, 2)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	if _, err := h.Dehumanize([]string{"a", "b"}); err == nil {
		t.Fatal("expected an eror, got nil")
	}
}

func TestHumanize(t *testing.T) {
	var hash = []byte{
		// b01010101 = 85
		85,

		// b10101010 = 170
		170,
	}

	words := []string{
		"a",
		"b",
	}

	expected := []string{
		"a",
		"b",
		"a",
		"b",
		"a",
		"b",
		"a",
		"b",
		"b",
		"a",
		"b",
		"a",
		"b",
		"a",
		"b",
		"a",
	}

	h, err := New(words, len(hash))
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	humanized, err := h.Humanize(hash)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	if len(humanized) != len(expected) {
		t.Fatalf("expected length %d, got: %d", len(expected), len(humanized))
	}

	for i := 0; i < len(expected); i++ {
		if humanized[i] != expected[i] {
			t.Fatalf("for %d expected %v but got %v", i, expected[i], humanized[i])
		}
	}
}

func TestDehumanizeInvalidWords(t *testing.T) {
	words := []string{
		"a",
		"b",
	}

	h, err := New(words, 5)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	humanized := createWords(h.NumberOfWords())
	_, err = h.Dehumanize(humanized)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func test(words []string, hash []byte, t *testing.T) {
	h, err := New(words, len(hash))
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	humanized, err := h.Humanize(hash)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	expectedNumberOfWords := h.NumberOfWords()
	if len(humanized) != expectedNumberOfWords {
		t.Fatalf("expected %d words, got: %d", expectedNumberOfWords, len(humanized))
	}

	dehumanized, err := h.Dehumanize(humanized)
	if err != nil {
		t.Fatalf("expected nil, got: %s", err)
	}

	expectedNumberOfBytes := h.NumberOfBytes()
	if len(dehumanized) != expectedNumberOfBytes {
		t.Fatalf("expected %d bytes, got: %d", expectedNumberOfBytes, len(dehumanized))
	}

	if !bytes.Equal(dehumanized, hash) {
		t.Fatalf("got %x expected %x", dehumanized, hash)
	}
}

func createWords(n int) []string {
	var rv []string
	for i := 0; i < n; i++ {
		rv = append(rv, fmt.Sprintf("%d", i))
	}
	return rv
}

func TestOneByte(t *testing.T) {
	const numberOfWords = 100
	for nWords := 2; nWords <= numberOfWords; nWords++ {
		fmt.Printf("words=%d\n", nWords)
		words := createWords(nWords)
		for i := 0; i < 256; i++ {
			hash := []byte{byte(i)}
			test(words, hash, t)
		}
	}
}

func TestTwoBytes(t *testing.T) {
	for n := 1; n <= 10; n++ {
		nWords := int(math.Pow(2, float64(n)))
		fmt.Printf("words=%d\n", nWords)
		words := createWords(nWords)
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				hash := []byte{byte(i), byte(j)}
				test(words, hash, t)
			}
		}
	}
}

func checkBits(b byte, expected []bool, t *testing.T) {
	for i := 0; i < 8; i++ {
		result := checkBit(b, i)
		if result != expected[i] {
			t.Fatalf("for %d expected %v but got %v", i, expected[i], result)
		}
	}
}

func TestCheckBit(t *testing.T) {
	// b01010101 = 85
	var b byte = 85

	expected := []bool{
		false,
		true,
		false,
		true,
		false,
		true,
		false,
		true,
		false,
		true,
	}

	checkBits(b, expected, t)
}

func TestSetBit(t *testing.T) {
	// b01010101 = 85
	var b byte = 85

	expected := []bool{
		false,
		true,
		true,
		true,
		false,
		true,
		false,
		true,
		false,
		true,
	}

	b = setBit(b, 2)

	checkBits(b, expected, t)
}

func TestClearBit(t *testing.T) {
	// b01010101 = 85
	var b byte = 85

	expected := []bool{
		false,
		false,
		false,
		true,
		false,
		true,
		false,
		true,
		false,
		true,
	}

	b = clearBit(b, 1)

	checkBits(b, expected, t)
}

// Package coupon can create formal coupon code with 4-character grouped parts.
package coupon

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	defaultNumParts = 3
)

var (
	defaultBadRegex = makeBadRegex(nil, false)
	symbols         = []byte("0123456789ABCDEFGHJKLMNPQRTUVWXY")
	alterGenerators = []regexReplacer{
		{regexp.MustCompile("[I1]"), "[I1]"},
		{regexp.MustCompile("[O0]"), "[O0]"},
		{regexp.MustCompile("[S5]"), "[S5]"},
		{regexp.MustCompile("[Z2]"), "[Z2]"},
	}
)

type regexReplacer struct {
	regex   *regexp.Regexp
	replace string
}

func (r *regexReplacer) ReplaceAll(original string) string {
	return r.regex.ReplaceAllString(original, r.replace)
}

func makeBadRegex(extraWords []string, replace bool) *regexp.Regexp {
	wordList := defaultBadWordList()
	if len(extraWords) > 0 {
		if replace {
			wordList = extraWords
		} else {
			wordList = append(wordList, extraWords...)
		}
	}

	for i := range wordList {
		for _, r := range alterGenerators {
			wordList[i] = r.ReplaceAll(wordList[i])
		}
	}

	return regexp.MustCompile(fmt.Sprintf("^(?:%s)$", strings.Join(wordList, "|")))
}

func defaultBadWordList() []string {
	words := []string{
		"SHPX", "PHAG", "JNAX", "JNAT", "CVFF", "PBPX", "FUVG", "GJNG", "GVGF", "SNEG", "URYY", "ZHSS", "QVPX", "XABO",
		"NEFR", "FUNT", "GBFF", "FYHG", "GHEQ", "FYNT", "PENC", "CBBC", "OHGG", "SRPX", "OBBO", "WVFZ", "WVMM", "CUNG",
	}

	for i := range words {
		for j := range words[i] {
			// rot13
			[]byte(words[i])[j] = (words[i][j]-'a'+13)%('z'-'a'+1) + 'a'
		}
	}

	return words
}

// Generator creates a coupon code by given members.
type Generator struct {
	// The number of parts. Must be a number in the range 1-6.
	// Default is 3.
	NumParts int

	// Seed for a coupon code. It hashed at first using sha1, and used as coupon index table.
	// If it is empty string then a randomized string is used.
	PlainText string

	// You can supply a regex for each part to not be generated.
	// If not provided then predefined regex is used.
	BadRegex *regexp.Regexp
}

func (g *Generator) Generate() string {
	plainText := g.PlainText
	if plainText == "" {
		plainText = g.randomPlainText()
	}

	badWords := g.BadRegex
	if badWords == nil {
		badWords = defaultBadRegex
	}

	numParts := g.NumParts
	if numParts < 1 {
		numParts = defaultNumParts
	} else if numParts > 6 {
		numParts = 6
	}

	hash := sha1.Sum([]byte(plainText))
	b := []byte(hex.EncodeToString(hash[:]))
	for i := range b {
		b[i] = b[i] % byte(len(symbols))
	}

	code := make([]string, 0, numParts)
	for len(code) < numParts {
		i := len(code) + 1
		str := []byte{symbols[b[0]], symbols[b[1]], symbols[b[2]]}
		b = b[3:]
		part := append(str, checkDigitAlg1(str, i))

		if badWords.Match(part) {
			continue
		} else if validWhenSwapped(part, i) {
			continue
		}

		code = append(code, string(part))
	}

	return strings.Join(code, "-")
}

func (g *Generator) randomPlainText() string {
	return fmt.Sprintf("%d%d", time.Now().Unix(), rand.Int63())
}

func checkDigitAlg1(data []byte, pos int) byte {
	check := pos
	for _, i := range []int{0, 1, 2} {
		for k := range symbols {
			if symbols[k] == data[i] {
				check = check*19 + k
				break
			}
		}
	}

	return symbols[check%len(symbols)]
}

func validWhenSwapped(orig []byte, pos int) bool {
	for _, code := range [][]byte{
		{orig[1], orig[0], orig[2], orig[3]},
		{orig[0], orig[2], orig[1], orig[3]},
		{orig[0], orig[1], orig[3], orig[2]},
	} {
		if string(code) == string(orig) {
			continue
		}

		if checkDigitAlg1(code[:3], pos) == code[3] {
			return true
		}
	}

	return false
}

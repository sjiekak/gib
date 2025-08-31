// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

// Package gib implements a gibberish string evaluator.
package gib

import (
	"math"
	"strings"
)

var (
	lowerCaseLetters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i",
		"j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w",
		"x", "y", "z"}
)

// allNgrams returns a list of all possible n-grams.
func allNgrams(n int) []string {
	if n == 0 {
		return []string{}
	} else if n == 1 {
		return lowerCaseLetters
	}

	newNgrams := make([]string, 0)
	for _, letter := range lowerCaseLetters {
		for _, ngram := range allNgrams(n - 1) {
			newNgrams = append(newNgrams, letter+ngram)
		}
	}
	return newNgrams
}

// ngramIDFValue computes scores using modified TF-IDF.
func ngramIDFValue(totalStrings, stringFreq, totalFreq, maxFreq float64) float64 {
	return math.Log2(totalStrings / (1. + stringFreq))
}

// highestIDF computes highest idf value in map of ngram frequencies.
func highestIDF(ngramFreq NGramScores) float64 {

	max := 0.
	for _, ngram := range ngramFreq {
		max = math.Max(max, ngram.IDF())
	}
	return max
}

// nGramValues computes n-gram statistics across a given corpus of strings.
func nGramValues(corpus []string, n int, reAdjust bool) NGramScores {
	var counts = make(map[string]int, n)
	var occurrences = NewNGramSet()
	var numStrings int

	for _, s := range corpus {
		s = strings.ToLower(s)
		numStrings++
		for _, ngram := range ngramsFromString(s, n) {
			occurrences.Add(ngram, s)
			counts[ngram]++
		}
	}

	keys := allNgrams(n)
	values := make([]Score, len(keys))

	generatedNGrams := NewNGramDict(keys, values)
	maxFreq := 0
	// computes max count and assign it as max frequency of ngram
	for _, k := range counts {
		maxFreq = int(math.Max(float64(k), float64(maxFreq)))
	}

	for ngram, strings := range occurrences.Set {
		stringFreq := len(strings)
		totalFreq := counts[ngram]
		score := ngramIDFValue(float64(numStrings), float64(stringFreq),
			float64(totalFreq), float64(maxFreq))
		generatedNGrams[ngram] = [3]float64{
			float64(stringFreq),
			float64(totalFreq),
			score,
		}
	}

	if reAdjust {
		maxIDF := math.Ceil(highestIDF(generatedNGrams))
		for ngram, value := range generatedNGrams {
			if value.IDF() == 0 {
				generatedNGrams[ngram] = [3]float64{
					0.,
					0.,
					maxIDF,
				}
			}
		}
	}

	return generatedNGrams
}

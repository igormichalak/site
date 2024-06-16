package main

import (
	"math"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/igormichalak/site/view"
)

func overlap(s1 []string, s2 []string) int {
	if s1 == nil || s2 == nil {
		return 0
	}
	total := 0
	set := make(map[string]struct{})
	for _, str := range s1 {
		set[str] = struct{}{}
	}
	for _, str := range s2 {
		if _, exists := set[str]; exists {
			total++
		}
	}
	return total
}

func distance(a string, b string) int {
	aLen := utf8.RuneCountInString(a)
	bLen := utf8.RuneCountInString(b)

	d := make([][]int, aLen+1)

	for i := 0; i <= aLen; i++ {
		d[i] = make([]int, bLen+1)
		for j := 0; j <= bLen; j++ {
			if i == 0 {
				d[i][j] = j
			} else if j == 0 {
				d[i][j] = i
			} else {
				d[i][j] = 0
			}
		}
	}

	for i := 1; i <= aLen; i++ {
		for j := 1; j <= bLen; j++ {
			var cost int
			if a[i-1] == b[j-1] {
				cost = 0
			} else {
				cost = 1
			}

			d[i][j] = min(
				d[i-1][j]+1,
				d[i][j-1]+1,
				d[i-1][j-1]+cost,
			)

			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+cost)
			}
		}
	}

	return d[aLen][bLen]
}

func cumulativeDistance(a, b []string) int {
	total := 0
	for _, aa := range a {
		smallest := math.MaxInt
		for _, bb := range b {
			smallest = min(smallest, distance(aa, bb))
		}
		total += smallest
	}
	return total
}

func words(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return strings.Fields(s)
}

func similaritySort(posts []view.PostSearchEntry, target string) {
	targetWords := words(target)

	slices.SortStableFunc(posts, func(a, b view.PostSearchEntry) int {
		return cumulativeDistance(targetWords, words(a.Title)) - cumulativeDistance(targetWords, words(b.Title))
	})
}

func filterByTags(posts []view.PostSearchEntry, tags []string) []view.PostSearchEntry {
	type entryWithTagCount struct {
		post  view.PostSearchEntry
		count int
	}

	var filtered []entryWithTagCount

	for _, post := range posts {
		var normalizedTags []string
		for _, tag := range post.Tags {
			normalizedTags = append(normalizedTags, view.NormalizeTag(tag))
		}

		matches := overlap(normalizedTags, tags)
		if matches == 0 {
			continue
		}
		filtered = append(filtered, entryWithTagCount{
			post:  post,
			count: matches,
		})
	}

	slices.SortStableFunc(filtered, func(a, b entryWithTagCount) int {
		return b.count - a.count
	})

	result := make([]view.PostSearchEntry, len(filtered))
	for i, entry := range filtered {
		result[i] = entry.post
	}

	return result
}

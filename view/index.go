package view

import (
	"fmt"
	"slices"
	"time"

	"github.com/a-h/templ"
)

type PostEntry struct {
	Component templ.Component
	Title     string
	Tags      []string
	CreatedAt time.Time
}

type PostSearchEntry struct {
	URL       string
	Title     string
	Tags      []string
	CreatedAt time.Time
}

var PostIndex = make(map[string]PostEntry)

var AllPostEntries []PostSearchEntry
var AllTags []string

func init() {
	for slug, entry := range PostIndex {
		AllPostEntries = append(AllPostEntries, PostSearchEntry{
			URL:       fmt.Sprintf("https://www.igormichalak.com/%s", slug),
			Title:     entry.Title,
			Tags:      entry.Tags,
			CreatedAt: entry.CreatedAt,
		})
		for _, tag := range entry.Tags {
			if !slices.Contains(AllTags, tag) {
				AllTags = append(AllTags, tag)
			}
		}
	}
	slices.SortFunc(AllPostEntries, func(a, b PostSearchEntry) int {
		return int(b.CreatedAt.Unix() - a.CreatedAt.Unix())
	})
}

package view

import (
	"time"

	"github.com/a-h/templ"
)

type PostEntry struct {
	Component   templ.Component
	Title       string
	Description string
	Tags        []string
	CreatedAt   time.Time
}

var PostIndex = make(map[string]PostEntry)

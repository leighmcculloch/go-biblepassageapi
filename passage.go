package biblepassageapi

import "strings"

type Passage struct {
	HTML         string
	TrackingCode string
	Copyright    string
}

func (p Passage) TimeToReadInMinutes() int {
	const readingWordsPerMin = 220
	text := p.HTML
	wordCount := strings.Count(text, " ")
	minutes := wordCount / readingWordsPerMin
	if minutes == 0 {
		minutes = 1
	}
	return minutes
}

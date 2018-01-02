package biblepassageapi

import "strings"

type Passage struct {
	Html         string
	TrackingCode string
	Copyright    string
}

func (p *Passage) TimeToReadInMinutes() int {
	const READING_WORDS_PER_MINUTE = 220
	text := p.Html
	wordCount := strings.Count(text, " ")
	wordsPerMinute := wordCount / READING_WORDS_PER_MINUTE
	if wordsPerMinute == 0 {
		wordsPerMinute = 1
	}
	return wordsPerMinute
}

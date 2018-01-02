package biblepassageapi

type Bible interface {
	Source() string
	NameShort() string
	NameCommon() string
	Name() string
	GetPassage(reference string) (*Passage, error)
}

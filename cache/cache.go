package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"4d63.com/biblepassageapi"
)

func Cache(b biblepassageapi.Bible, cacheFolderPath string) biblepassageapi.Bible {
	return &bible{Bible: b, cacheFolderPath: cacheFolderPath}
}

type bible struct {
	biblepassageapi.Bible
	cacheFolderPath string
}

func (b *bible) GetPassage(reference string) (biblepassageapi.Passage, error) {
	return getBiblePassageWithCache(b.cacheFolderPath, b.Bible, reference)
}

func getBiblePassageWithCache(cacheFolderPath string, b biblepassageapi.Bible, reference string) (biblepassageapi.Passage, error) {
	p, err := loadBiblePassage(cacheFolderPath, b, reference)
	if err != nil {
		p, err := b.GetPassage(reference)
		if err != nil {
			return biblepassageapi.Passage{}, err
		}
		err = saveBiblePassage(cacheFolderPath, b, reference, p)
		if err != nil {
			return biblepassageapi.Passage{}, err
		}
		return p, nil
	}
	return p, nil
}

func loadBiblePassage(cacheFolderPath string, b biblepassageapi.Bible, reference string) (biblepassageapi.Passage, error) {
	f, err := os.Open(getCacheFilePath(cacheFolderPath, b, reference))
	if err != nil {
		return biblepassageapi.Passage{}, err
	}
	defer f.Close()

	var p biblepassageapi.Passage
	dec := gob.NewDecoder(f)
	err = dec.Decode(&p)
	if err != nil {
		return biblepassageapi.Passage{}, err
	}

	return p, nil
}

func saveBiblePassage(cacheFolderPath string, b biblepassageapi.Bible, reference string, p biblepassageapi.Passage) error {
	err := os.MkdirAll(getCacheFileDir(cacheFolderPath, b, reference), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(getCacheFilePath(cacheFolderPath, b, reference))
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	enc.Encode(p)
	if err != nil {
		return err
	}

	return nil
}

func getCacheFileDir(cacheFolderPath string, b biblepassageapi.Bible, reference string) string {
	return filepath.Join(cacheFolderPath, b.Source(), b.NameShort())
}

func getCacheFileName(b biblepassageapi.Bible, reference string) string {
	return fmt.Sprintf("%s.biblepassage", reference)
}

func getCacheFilePath(cacheFolderPath string, b biblepassageapi.Bible, reference string) string {
	return filepath.Join(getCacheFileDir(cacheFolderPath, b, reference), getCacheFileName(b, reference))
}

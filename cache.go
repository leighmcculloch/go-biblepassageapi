package biblepassageapi

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

type bibleCache struct {
	Bible
	cacheFolderPath string
}

func Cache(b Bible, cacheFolderPath string) Bible {
	return &bibleCache{Bible: b, cacheFolderPath: cacheFolderPath}
}

func (bc *bibleCache) GetPassage(reference string) (*Passage, error) {
	return getBiblePassageWithCache(bc.cacheFolderPath, bc.Bible, reference)
}

func getBiblePassageWithCache(cacheFolderPath string, b Bible, reference string) (*Passage, error) {
	p, err := loadBiblePassage(cacheFolderPath, b, reference)
	if err != nil {
		p, err := b.GetPassage(reference)
		if err != nil {
			return nil, err
		}
		err = saveBiblePassage(cacheFolderPath, b, reference, p)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return p, nil
}

func loadBiblePassage(cacheFolderPath string, b Bible, reference string) (*Passage, error) {
	f, err := os.Open(getCacheFilePath(cacheFolderPath, b, reference))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bp Passage
	dec := gob.NewDecoder(f)
	err = dec.Decode(&bp)
	if err != nil {
		return nil, err
	}

	return &bp, nil
}

func saveBiblePassage(cacheFolderPath string, b Bible, reference string, p *Passage) error {
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

func getCacheFileDir(cacheFolderPath string, b Bible, reference string) string {
	return filepath.Join(cacheFolderPath, b.Source(), b.NameShort())
}

func getCacheFileName(b Bible, reference string) string {
	return fmt.Sprintf("%s.biblepassage", reference)
}

func getCacheFilePath(cacheFolderPath string, b Bible, reference string) string {
	return filepath.Join(getCacheFileDir(cacheFolderPath, b, reference), getCacheFileName(b, reference))
}

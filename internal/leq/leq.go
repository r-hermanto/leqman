package leq

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Collection struct {
	Name     string
	FilePath string
}

func GetCollections(path string) []Collection {
	collections := []Collection{}

	de, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, itm := range de {
		if itm.IsDir() {
			continue
		}

		ext := filepath.Ext(itm.Name())
		if ext != ".json" {
			continue
		}

		leqName := strings.TrimSuffix(filepath.Base(itm.Name()), ext)
		collections = append(collections, Collection{
			Name:     leqName,
			FilePath: path + itm.Name(),
		})
	}

	return collections
}

type LeqConfig struct {
	URL    string          `json:"url"`
	Method string          `json:"method"`
	Header json.RawMessage `json:"header"`
	Body   json.RawMessage `json:"body"`
}

func NewLeqConfig(l *LeqConfig) {

}

func (l *LeqConfig) upsertRequest() {

}

func (l *LeqConfig) deleteRequest() {

}

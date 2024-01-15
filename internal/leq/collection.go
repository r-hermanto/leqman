package leq

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Collection struct {
	Title    string
	Path     string
	IsDir    bool
	Children []*Collection
}

func GetCollections() []*Collection {
	path := getLeqConfigDir()

	var generateCollections func(string) []*Collection
	generateCollections = func(path string) []*Collection {
		collections := []*Collection{}

		dirEntry, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, itm := range dirEntry {
			fullPath := filepath.Join(path, itm.Name())
			collection := &Collection{
				Title: getTitle(itm.Name()),
				Path:  fullPath,
				IsDir: itm.IsDir(),
			}

			if itm.IsDir() {
				collection.Children = generateCollections(fullPath)
			}

			collections = append(collections, collection)
		}

		return collections
	}

	return generateCollections(path)
}

func getTitle(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}

	if ext != ".json" {
		return ""
	}

	return strings.TrimSuffix(filepath.Base(name), ext)
}

func getLeqConfigDir() string {
	// TODO: allow user to define custom path in xdg config
	leqmanDir := os.Getenv("HOME") + "/leqman/"
	os.Mkdir(leqmanDir, os.ModePerm)
	return leqmanDir
}

func GetRequest(filePath string) LeqConfig {
	if _, err := os.Stat(filePath); errors.Is(err, fs.ErrNotExist) {
		log.Fatal("Request not exists.", err)
	}

	fileByte, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var legConfig LeqConfig
	err = json.Unmarshal(fileByte, &legConfig)
	if err != nil {
		log.Fatal(err)
	}

	return legConfig
}

func DeleteRequest(name string) {
	path := getLeqConfigDir() + name + ".json"
	if err := os.Remove(path); err != nil {
		log.Fatal(err)
	}
}

func CreateNewRequest(name string, l LeqConfig) {
	path := getLeqConfigDir() + name + ".json"
	if _, err := os.Stat(path); err == nil {
		log.Fatal("Request already exists.", err)
	}

	writeRequestToFile(path, l)
}

func UpdateRequest(name string, l LeqConfig) {
	path := getLeqConfigDir() + name + ".json"
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		log.Fatal("Request not exists.", err)
	}

	writeRequestToFile(path, l)
}

func writeRequestToFile(path string, l LeqConfig) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := json.Marshal(l)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

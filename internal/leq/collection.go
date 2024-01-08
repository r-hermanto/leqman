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
	Name     string
	FilePath string
}

func GetCollections() []Collection {
	collections := []Collection{}

	path := getLeqConfigDir()
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

func getLeqConfigDir() string {
	// TODO: allow user to define custom path in xdg config
	leqmanDir := os.Getenv("HOME") + "/leqman/"
	os.Mkdir(leqmanDir, os.ModePerm)
	return leqmanDir
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

package main

import (
	"fmt"
	"os"

	"github.com/r-hermanto/leqman/internal/leq"
)

func main() {
	dir := getLeqConfigDir()
	collections := leq.GetCollections(dir)
	fmt.Printf("collections: %v\n", collections)
}

func getLeqConfigDir() string {
	// TODO: allow user to define custom path in xdg config
	leqmanDir := os.Getenv("HOME") + "/leqman/"
	os.Mkdir(leqmanDir, os.ModePerm)
	return leqmanDir
}

package main

import (
	"fmt"

	"github.com/r-hermanto/leqman/internal/leq"
)

func main() {
	collections := leq.GetCollections()
	fmt.Printf("collections: %v\n", collections)

	config := leq.LeqConfig{
		URL:    "locahost:8000/home1",
		Body:   nil,
		Header: nil,
		Method: "GET",
	}
	leq.UpdateRequest("homepage", config)

	collections = leq.GetCollections()
	fmt.Printf("collections: %v\n", collections)

	leq.DeleteRequest("homepage")

	collections = leq.GetCollections()
	fmt.Printf("collections: %v\n", collections)
}

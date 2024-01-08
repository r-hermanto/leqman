package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/r-hermanto/leqman/internal/leq"
)

func main() {
	collections := leq.GetCollections()
	fmt.Printf("collections: %v\n", collections)

	body := map[string]interface{}{
		"body1": "value body",
		"body2": 123,
		"body3": map[string]interface{}{
			"inner1": 123,
			"inner2": "inner value 2",
		},
	}

	bbody, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	config := leq.LeqConfig{
		URL:  "localhost:8000/home",
		Body: bbody,
		Header: map[string]string{
			"header1": "value1",
			"header2": "value2",
		},
		Method: "GET",
	}
	leq.UpdateRequest("homepage", config)

	api := leq.GetRequest("homepage")
	api.Execute()
}

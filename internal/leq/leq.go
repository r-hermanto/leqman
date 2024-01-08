package leq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
)

type LeqConfig struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   json.RawMessage   `json:"body"`
}

func (l *LeqConfig) Execute() {
	url := l.URL
	if url[:4] != "http" {
		url = "http://" + url
	}

	body := &bytes.Buffer{}
	if allowBody(l.Method) {
		body = bytes.NewBuffer(l.Body)
	}

	req, err := http.NewRequestWithContext(context.TODO(), l.Method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range l.Header {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var pretty bytes.Buffer
	json.Indent(&pretty, b, "", "    ")
	fmt.Printf("res.Header: %v\n", resp.Header)
	fmt.Println(pretty.String())
}

func allowBody(httpMethod string) bool {
	allowed := []string{http.MethodPost, http.MethodPut, http.MethodPatch}
	return slices.Contains(allowed, httpMethod)
}

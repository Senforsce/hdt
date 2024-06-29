package crawler

import (
	"fmt"
	"os"
	"regexp"
)

type OntologyFileClassifier struct {
	Prefixes []string
	ClassifierRunner
}

func NewClassifierForCrawler(prefixes []string) *OntologyFileClassifier {
	f := &OntologyFileClassifier{
		Prefixes: prefixes,
	}

	return f
}

func (h OntologyFileClassifier) Name() string {
	return "O8"
}

func (h OntologyFileClassifier) Run(path string, ContextStore map[string]map[string]string, key string) {
	b, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	if ContextStore[key] == nil {
		ContextStore[key] = make(map[string]string)
	}

	for _, ns := range h.Prefixes {
		handleO8Namespace(
			fmt.Sprintf("%s:[A-Za-z0-9-]+", ns),
			b, ContextStore[key])
	}
}

// TODO: use io.Reader instead
func handleO8Namespace(reg string, b []byte, s map[string]string) {
	regExp, _ := regexp.Compile(reg)
	matches := regExp.FindAll(b, -1)
	if len(matches) > 0 {
		mutateWithMatches(s, matches)
	}
}

// TODO: use io.Reader instead
func mutateWithMatches(s map[string]string, matches [][]byte) {
	for _, v := range matches {
		s[string(v)] = string(v)
	}
}

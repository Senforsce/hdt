package crawler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PreviewFileClassifier struct {
	ClassifierRunner
	ImportPrefix string
}

func NewClassifierForPreviews(pr string) *PreviewFileClassifier {
	f := &PreviewFileClassifier{
		ImportPrefix: pr,
	}

	return f
}

func (h PreviewFileClassifier) Name() string {
	return "Preview"
}

func (h PreviewFileClassifier) Run(path string, ContextStore map[string]map[string]string, key string) {
	dir := filepath.Dir(path)
	base := filepath.Base(dir)
	foundPreview := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		regExp, _ := regexp.Compile("_preview.go$")
		matches := regExp.FindAll([]byte(path), -1)

		if !info.IsDir() && len(matches) > 0 {
			foundPreview = true
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	remote := h.ImportPrefix + "owl/" + base
	importpath := strings.Replace(remote, h.ImportPrefix, "/", 1)
	route := strings.Replace(importpath, "owl", "preview", 1)
	if ContextStore[key] == nil {
		ContextStore[key] = make(map[string]string)
	}
	if foundPreview {
		ContextStore[key][base] = route + "::" + remote

	}
}

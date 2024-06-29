package crawler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SpiderConfig struct {
	RelativeFilePath    string
	Prefixes            []string
	LookupFileExtension string
	Classifiers         []ClassifierRunner
}

type ClassifierRunner interface {
	Run(path string, store map[string]map[string]string, key string)
	Name() string
}

// walks a defined directory and finds a file extension
// looks-up a set of ontology namespaces and returns a map of files to be treated
func FileOntologySpider(cfg *SpiderConfig) map[string]map[string]string {
	var root string
	// Replace 'filename' with the actual filename you want to use
	dir, r := os.Getwd()
	if r != nil {
		fmt.Println("Cannot get current directory")
		return map[string]map[string]string{}
	}
	// Small security enforcement, prevent the user to go further than the parent directory
	// Only the line below should allow ..
	boundary := "/.."
	root = dir + boundary + strings.ReplaceAll(cfg.RelativeFilePath, "..", "")
	var files []string
	var ContextStore = make(map[string]map[string]string)

	// hdt := "../main.hdt"
	// base := "http://senforsce.com/o8/Brain"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		ext := fmt.Sprintf(".%s", cfg.LookupFileExtension)
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
			for _, c := range cfg.Classifiers {
				c.Run(path, ContextStore, c.Name())
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		log.Println(file)
	}

	return ContextStore
}

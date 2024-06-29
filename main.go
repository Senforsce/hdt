package main

/*
#cgo LDFLAGS: -L${SRCDIR}/libhdt/.libs -lhdt -lstdc++
#include <stdlib.h>
#include "hdtwrapper.h"
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"unsafe"

	"github.com/fatih/color"

	"github.com/senforsce/hdt/crawler"
)

type Query struct {
	Q string
}

type TripleResponse struct {
	Success bool
	Triples []string
}

func main() {
	// Replace 'filename' with the actual filename you want to use
	filename := "../main.env.tmpl.ttl"
	out := "../main.hdt"
	base := "http://senforsce.com/o8/brain"
	// Convert Go string to C string
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cOut := C.CString(out)
	defer C.free(unsafe.Pointer(cOut))

	cBase := C.CString(base)
	defer C.free(unsafe.Pointer(cBase))

	result := C.generateHDTWrapper(cFilename, cBase, cOut)
	// result := 0 //C.generateHDTWrapper(cFilename, cBase, cOut)
	// Check the result
	if result != 0 {
		fmt.Println("Error generating HDT")
	} else {
		fmt.Println("HDT generated successfully")
		autoInjectOntologyReferences()
	}
}

func autoInjectOntologyReferences() error {
	// TODO: extract prefixes from .ttl file
	prefixes := []string{"i18n", "path", "se", "sc", "sen", "hx", "id", "ids", "cls", "clss", "role",
		"tid"}
	namespaceList := map[string]string{
		"i18n:": "http://senforsce.com/o8/brain/i18n/",
		"path:": "http://senforsce.com/o8/brain/path/",
		"se:":   "http://senforsce.com/o8/brain/SecureEnvironment/",
		"sc:":   "http://senforsce.com/o8/brain/SecureConfig/",
		"hx:":   "http://senforsce.com/o8/brain/HTMXAttribute/",
		"id:":   "http://senforsce.com/o8/brain/Identifier/",
		"ids:":  "http://senforsce.com/o8/brain/IdSelector/",
		"cls:":  "http://senforsce.com/o8/brain/Class/",
		"clss:": "http://senforsce.com/o8/brain/ClassSelector/",
		"role:": "http://senforsce.com/o8/brain/AriaRole/",
		"tid:":  "http://senforsce.com/o8/brain/TestId/",
		"sen:":  "http://senforsce.com/o8/brain/",
	}
	// todo: read these values from config
	relativeFilePath := "/tndrf1sh"
	InputFile := "../main.hdt"
	fileExtensionForLookingUpPrefixes := "t1"
	classifiers := []crawler.ClassifierRunner{}

	// todo decouple logic from generating hdt to spidering
	classifiers = append(classifiers, *crawler.NewClassifierForCrawler(prefixes))
	classifiers = append(classifiers, *crawler.NewClassifierForPreviews("github.com/senforsce/tndrf1sh/web/"))

	spiderConf := &crawler.SpiderConfig{
		Prefixes:            prefixes,
		RelativeFilePath:    relativeFilePath,
		LookupFileExtension: fileExtensionForLookingUpPrefixes,
		Classifiers:         classifiers,
	}
	store := crawler.FileOntologySpider(spiderConf)
	// var x string
	var s string
	var errList = make(map[string]string)
	var errCount = make(map[string]int)
	var injectable = make(map[string]string)
	var allErr int = 0
	var lookup = map[string]string{
		"i18n": "http://senforsce.com/o8/brain/enTranslation",
		"se":   "http://senforsce.com/o8/brain/hasUrl",
		"sc":   "http://senforsce.com/o8/brain/hasValue",
		"path": "http://senforsce.com/o8/brain/hasUrl",
		"hx":   "http://senforsce.com/o8/brain/val",
		"id":   "http://senforsce.com/o8/brain/val",
		"ids":  "http://senforsce.com/o8/brain/val",
		"cls":  "http://senforsce.com/o8/brain/val",
		"clss": "http://senforsce.com/o8/brain/val",
		"role": "http://senforsce.com/o8/brain/val",
		"tid":  "http://senforsce.com/o8/brain/val",

		"sen": "http://senforsce.com/o8/brain/enTranslation",
	}
	for key := range store["O8"] {
		s = key
		errCount[key] = 0
		errList[key] = ""
		for k, v := range namespaceList {
			s = strings.Replace(s, k, v, -1)
			fmt.Printf("replace:( %s ) %s -> %s = %s\n", key, k, v, s)
		}
		full := strings.Split(key, ":")
		prefix := full[0]
		word := full[1]

		p := lookup[prefix]
		o := ""
		cS := C.CString(s)
		cP := C.CString(p)
		cInputFile := C.CString(InputFile)
		defer C.free(unsafe.Pointer(cInputFile))
		cO := C.CString(o)
		defer C.free(unsafe.Pointer(cS))
		defer C.free(unsafe.Pointer(cP))
		defer C.free(unsafe.Pointer(cO))

		// Allocate a buffer for the result
		bufferSize := 1024 // Adjust this size according to your requirements
		resultBuffer := C.CString(string(make([]byte, bufferSize)))
		defer C.free(unsafe.Pointer(resultBuffer))

		ts := C.searchWrapper(cInputFile, cS, cP, cO, resultBuffer, C.int(bufferSize))
		notNeutral := (prefix != "ids" && prefix != "id" && prefix != "cls" && prefix != "clss")
		if ts != 0 {

			if notNeutral {
				msg := fmt.Sprintf("Error: ====> %s in %s, %s", key, InputFile, s)
				color.Red(msg)
				errCount[key] = errCount[key] + 1
				errList[key] = msg
				allErr = allErr + 1
			} else {
				color.Blue(fmt.Sprintf("neutral: ====> %s \n", key))
				var selector = ""
				if prefix == "clss" {
					selector = "."
				}
				if prefix == "ids" {
					selector = "#"
				}

				injectable[key] = `"` + selector + word + `"`
			}

		} else {
			// if ts != 0 {
			color.Green(fmt.Sprintf("Success: ====> %s \n", key))
			goString := C.GoString(resultBuffer)

			injectable[key] = goString
		}

	}

	msg := fmt.Sprintf("\n\n\nErrors: ====> %d", allErr)
	log.Println(msg)

	for kk := range store["O8"] {
		if errCount[kk] > 0 {
			bad := fmt.Sprintf("Errors: ==NOT_FOUND===> %s", kk)
			fmt.Println(bad)
		} else {
			// add the value to a list of values to be generated
			ok := fmt.Sprintf("Success: =====%s====> %s", InputFile, kk)
			color.Green(ok)
		}

	}
	injectOntologyAsFile(injectable, "../tndrf1sh/web/inject.go")
	injectHandlersFromHTMXMentions(store["Preview"], "../tndrf1sh/web/inject_previews.go")
	msg2 := fmt.Sprintf("\n\n\nErrors: ====> %d", allErr)
	color.Red(msg2)

	return nil
}

func injectOntologyAsFile(values map[string]string, filepath string) {
	text := `package main

import (
	"github.com/senforsce/tndr0cean/router"
)

func WithO8(app *router.Nw) error {
	{{ range $k, $v := . }}app.S("{{$k}}", {{$v}})
	{{end}}
	//temporary TODO inject from .go file
	app.S("path:SparQlCurrentServer", "/sparql/current-server")
	app.S("path:SparQlCurrentUser", "/sparql/current-user")
	app.S("path:SparQlCurrentIntegrations", "/sparql/current-integrations")
	app.S("path:SparQlCurrentComponent", "/sparql/current-component")

	return nil
}
	`

	fmt.Printf("%v", values)
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error")
	}

	tmpl, err := template.New("test").Parse(text)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, values)
	if err != nil {
		panic(err)
	}

	f.Close()

}

type Injectable struct {
	Remote      string
	Route       string
	HandlerFunc string
}

func injectHandlersFromHTMXMentions(values map[string]string, filepath string) {
	text := `package main

import (
	"github.com/senforsce/tndr0cean/router"
	{{ range $k, $v := . }}"{{$v.Remote}}" {{end}}
)

func WithHXMXComponents(app *router.Tndr0cean) error {
	{{ range $k, $v := . }}app.Get("{{$v.Route}}", {{$v.HandlerFunc}}) 
	{{end}}

	return nil
}
	`

	fmt.Printf("%v", values)
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error")
	}

	toUse := []Injectable{}

	for ks, vs := range values {
		ps := strings.Split(vs, "::")
		in := Injectable{
			Route:       ps[0],
			Remote:      ps[1],
			HandlerFunc: ks + ".Preview",
		}
		toUse = append(toUse, in)
	}

	tmpl, err := template.New("test2").Parse(text)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, toUse)
	if err != nil {
		panic(err)
	}

	f.Close()

}

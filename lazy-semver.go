package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io/ioutil"
	"strings"
)

func main() {
	filePath, xPathPattern, baseVersionRegex := handleFlags()
	baseVersion := getBaseVersion(*filePath, *xPathPattern) // base version is major.minor version and looks like 1.0

	fmt.Printf("Args were filePath: %s, %s, %s\n", *filePath, *xPathPattern, *baseVersionRegex)
	fmt.Printf("Parsed base version is: %s\n", baseVersion)
}

func getBaseVersion(filePath string, xPathPattern string) string {
	// read file
	fileContentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error: Could not read file \"%s\"\n", filePath))
	}
	fileContentString := string(fileContentBytes)

	// process
	if xPathPattern == "" {
		return fileContentString
	} else if strings.HasSuffix(strings.ToLower(filePath), ".xml") {
		doc, err := xmlquery.Parse(strings.NewReader(fileContentString))
		if err != nil {
			panic(err)
		}
		baseVersionNode := xmlquery.FindOne(doc, xPathPattern)
		return baseVersionNode.InnerText()
	} else {
		panic("You passed neither an xml file, nor a plaintext with an empty argument xPathPattern. " +
			"These are to only types currently supported.")
	}
}

func handleFlags() (*string, *string, *string) {
	filePath := flag.String("filePath", "no filePath provided", "filePath to a filename "+
		"including the extensions, of which the base version should be read")
	xPathPattern := flag.String("xPathPattern", "", "XPath Pattern to select the version. If "+
		"empty, the whole content of the file will be used as base version string.")
	baseVersionRegex := flag.String("baseVersionRegex", "\\d+.\\d+",
		"the regex used to parse the base version (that is the major.minor version part)")
	flag.Parse()
	return filePath, xPathPattern, baseVersionRegex
}

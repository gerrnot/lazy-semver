package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"io/ioutil"
	"strings"
)

func main() {
	filePath, xPathPattern, baseVersionRegex := handleFlags()
	baseVersion := getBaseVersion(*filePath, *xPathPattern)

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
	} else {
		expr, err := xpath.Compile(xPathPattern)
		if err != nil {
			panic(err)
		}
		doc, err := xmlquery.Parse(strings.NewReader(fileContentString))
		if err != nil {
			panic(err)
		}
		baseVersionNode := expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
		fmt.Printf("baseVersionString=%s", baseVersionNode.Current().Value())
	}
	panic("Error: ")
	return "" /* This line will not be executed, but the go compiler complains otherwise */
}

func handleFlags() (*string, *string, *string) {
	filePath := flag.String("filePath", "no filePath provided", "filePath to a filename including the "+
		"extensions, of which the base version should be read")
	xPathPattern := flag.String("xPathPattern", "", "XPath Pattern to select the version. If "+
		"empty, the whole content of the file will be used as base version string.")
	baseVersionRegex := flag.String("baseVersionRegex", "\\d+.\\d+", "the regex used to parse the base version (that is the major.minor version part)")
	flag.Parse()
	return filePath, xPathPattern, baseVersionRegex
}

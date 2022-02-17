package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	filePath, xPathPattern, baseVersionRegex, timestampRFC3339 := handleFlags()
	baseVersion := getBaseVersion(filePath, xPathPattern, baseVersionRegex) /* major.minor part of SemVer */
	commitCount := getCommitCount(*filePath)                                /* patch  part of SemVer */
	timestampRFC3339String := getTimestampRFC3339String(*timestampRFC3339)  /* +build  part of SemVer */
	calculatedVersion := fmt.Sprintf("%s.%d%s", baseVersion, commitCount, timestampRFC3339String)
	fmt.Print(calculatedVersion) /* major.minor.patch[+UTC timestamp] */
}

func getTimestampRFC3339String(timestampRFC3339 bool) string {
	if timestampRFC3339 {
		t := time.Now()
		return "+" + t.Format(time.RFC3339)
	} else {
		return ""
	}
}

func getBaseVersion(filePath *string, xPathPattern *string, baseVersionRegex *string) string {
	// rawVersionString looks like 1.0.0-SNAPSHOT here
	rawVersionString := getOriginalVersionStringFromFile(*filePath, *xPathPattern)
	regex, err := regexp.Compile(*baseVersionRegex)
	if err != nil {
		panic(err)
	}
	baseVersion := regex.FindString(rawVersionString)
	// baseVersion looks like 1.0 here
	return baseVersion
}

func getOriginalVersionStringFromFile(filePath string, xPathPattern string) string {
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

func handleFlags() (*string, *string, *string, *bool) {
	filePath := flag.String("filePath", "no filePath provided", "filePath to a filename "+
		"including the extensions, of which the base version should be read")
	xPathPattern := flag.String("xPathPattern", "", "[optional] Only required when filePath"+
		"refers to a file that needs special parsing (e.g. to read base version from xml/json file). "+
		"XPath Pattern to select the version. If empty, the whole content of the file will be used as base "+
		"version string.")
	baseVersionRegex := flag.String("baseVersionRegex", "\\d+.\\d+",
		"[optional] the regex used to parse the base version (that is the major.minor version part)")
	timestampRFC3339 := flag.Bool("timestampRFC3339", false, "[optional] Attaches an ISO-1806 "+
		"timestamp as SemVer build reference. "+
		"Example of calculated version: 0.0.1+2022-02-17T11:12:27+01:00"+
		"The timestamp format is described here: "+
		"https://www.ietf.org/rfc/rfc3339.txt. Go defines a time format string with that name!")
	flag.Parse()
	return filePath, xPathPattern, baseVersionRegex, timestampRFC3339
}

func getCommitCount(filePath string) int {
	dirPath := filepath.Dir(filePath)
	repoPath := findGitRootRecursive(dirPath)
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		panic(err)
	}
	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}
	var cCount int
	err = cIter.ForEach(func(c *object.Commit) error {
		cCount++
		return nil
	})
	return cCount
}

func findGitRootRecursive(basePath string) string {
	if basePath == "" {
		panic("Could not find a .git directory! All parent folders of given parameter filePath were searched")
	}
	gitDirPath := filepath.Join(basePath, ".git")
	_, err := os.Stat(gitDirPath)
	if err != nil {
		return findGitRootRecursive(filepath.Dir(basePath))
	}
	return basePath
}

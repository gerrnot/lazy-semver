package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

var lazySemVerBinPath, testResourcesBasePath = getTestPaths()

func getTestPaths() (string, string) {
	currentPathFull, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(currentPathFull, "lazy-semver"),
		filepath.Join(currentPathFull, "test", "resources")
}

func Test_scenarioMavenWithTimestampAndResultFile(t *testing.T) {
	const resultFilePath = "/tmp/lazy-semver-test"
	cmd := exec.Command(
		lazySemVerBinPath,
		"--inputFilePath", filepath.Join(testResourcesBasePath, "pom.xml"),
		"--xPathPattern", "/project/version",
		"--timestampRFC3339",
		"--resultFilePath", resultFilePath,
	)
	cmd.Dir = testResourcesBasePath
	var stdOutBuffer bytes.Buffer
	var stdErrBuffer bytes.Buffer
	cmd.Stdout = &stdOutBuffer
	cmd.Stderr = &stdErrBuffer
	log.Infof("Executing Command: %s", cmd.String())
	errR := cmd.Run()
	stdout := stdOutBuffer.String()
	stderr := stdErrBuffer.String()
	// the last bit .* seems dependent on the system, when run in docker will yield Z (which stands for UTC timezone)
	// on my local machine it produced +01:00 (which also means UTC timezone, but obviously in a different format)
	expectedVersion, err := regexp.Compile("1.1.\\d+\\+\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*")
	if err != nil {
		panic(err)
	}
	if !expectedVersion.MatchString(stdout) {
		t.Errorf("Expected version to match regex '%s', but got '%s', which did not match. Stderr was: '%s'",
			expectedVersion, stdout, stderr)
	}
	if errR != nil {
		panic(err)
	}
	resultFileContent, err := os.ReadFile(resultFilePath)
	resultFileContentString := string(resultFileContent)
	if err != nil {
		panic(err)
	}
	if !expectedVersion.MatchString(resultFileContentString) {
		panic(fmt.Sprintf("Expected file %s to match %s, but got %s, which did not match.",
			resultFilePath, expectedVersion, resultFileContentString))
	}
}

func Test_scenarioPlainFile(t *testing.T) {
	cmd := exec.Command(
		lazySemVerBinPath,
		"--inputFilePath", filepath.Join(testResourcesBasePath, "version.txt"),
	)
	cmd.Dir = testResourcesBasePath
	var stdOutBuffer bytes.Buffer
	var stdErrBuffer bytes.Buffer
	cmd.Stdout = &stdOutBuffer
	cmd.Stderr = &stdErrBuffer
	log.Infof("Executing Command: %s", cmd.String())
	errR := cmd.Run()
	stdout := stdOutBuffer.String()
	stderr := stdErrBuffer.String()
	expectedVersion, err := regexp.Compile("3.100.\\d+")
	if err != nil {
		panic(err)
	}
	if !expectedVersion.MatchString(stdout) {
		t.Errorf("Expected version to match regex '%s', but got '%s' which did not match. Stderr was: '%s'",
			expectedVersion, stdout, stderr)
	}
	if errR != nil {
		panic(err)
	}
}

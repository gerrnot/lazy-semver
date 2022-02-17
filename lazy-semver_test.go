package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
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

func Test_scenarioMaven(t *testing.T) {
	cmd := exec.Command(
		lazySemVerBinPath,
		"--filePath", filepath.Join(testResourcesBasePath, "pom.xml"),
		"--xPathPattern", "/project/version",
	)
	cmd.Dir = testResourcesBasePath
	var stdOutBuffer bytes.Buffer
	var stdErrBuffer bytes.Buffer
	cmd.Stdout = &stdOutBuffer
	cmd.Stderr = &stdErrBuffer
	log.Infof("Executing Command: %s", cmd.String())
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	stdout := stdOutBuffer.String()
	stderr := stdErrBuffer.String()
	expectedVersion := "1.1.2"
	if stdout != expectedVersion {
		t.Errorf("Expected version '%s', but got '%s'. Stderr was: '%s'", expectedVersion, stdout, stderr)
	}
}

func Test_scenarioPlainFile(t *testing.T) {
	cmd := exec.Command(
		lazySemVerBinPath,
		"--filePath", filepath.Join(testResourcesBasePath, "version.txt"),
	)
	cmd.Dir = testResourcesBasePath
	var stdOutBuffer bytes.Buffer
	var stdErrBuffer bytes.Buffer
	cmd.Stdout = &stdOutBuffer
	cmd.Stderr = &stdErrBuffer
	log.Infof("Executing Command: %s", cmd.String())
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	stdout := stdOutBuffer.String()
	stderr := stdErrBuffer.String()
	expectedVersion := "3.100.2"
	if stdout != expectedVersion {
		t.Errorf("Expected version '%s', but got '%s'. Stderr was: '%s'", expectedVersion, stdout, stderr)
	}
}

package worker

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileFilterWhenFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "")
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	fileInfo, _ := os.Stat(tempFile.Name())

	result := FileFilter(tempFile.Name(), fileInfo)

	if !result {
		t.Error("filter result should be true")
	}
}

func TestFileFilterWhenDir(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)
	fileInfo, _ := os.Stat(tempDir)

	result := FileFilter(tempDir, fileInfo)

	if result {
		t.Error("filter result should be false")
	}
}

func TestDirFilterWhenFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile("", "")
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	fileInfo, _ := os.Stat(tempFile.Name())

	result := DirFilter(tempFile.Name(), fileInfo)

	if result {
		t.Error("filter result should be false")
	}
}

func TestDirFilterWhenDir(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)
	fileInfo, _ := os.Stat(tempDir)

	result := DirFilter(tempDir, fileInfo)

	if !result {
		t.Error("filter result should be true")
	}
}

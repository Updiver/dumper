package dumper

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var testRepository = "https://github.com/Updiver/test-repository.git"

// dummy test
func TestNew(t *testing.T) {
	dumper := New()

	if dumper == nil {
		t.Errorf("expect to have proper dump instance")
	}
}

func TestDumpRepository(t *testing.T) {
	var (
		destinationDir = "repository-clone-example"
	)

	dir := os.TempDir()
	fmt.Printf("temp dir: %s\n", dir)
	fmt.Printf("desitnation dir: %s\n", destinationDir)
	tempDir := path.Join(filepath.Clean(dir), destinationDir)
	defer os.RemoveAll(tempDir)

	dumper := New()
	err := dumper.DumpRepository(testRepository, tempDir, "", "")
	if err != nil {
		t.Errorf("dump repository: %v", err)
	}

	fileContent, err := os.Open(path.Join(tempDir, "test-regular-file.txt"))
	if err != nil {
		t.Errorf("open file: %v", err)
	}

	txt, err := io.ReadAll(fileContent)
	if err != nil {
		t.Errorf("read file content: %v", err)
	}

	fmt.Println(string(txt))
	if string(txt) != "Test regular file content" {
		t.Errorf("expect to have proper file content")
	}
}

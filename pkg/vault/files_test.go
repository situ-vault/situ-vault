package vault

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maybeFromFiles(t *testing.T) {
	f1 := testFile(t, "test-file-1-", "test-pw")
	f2 := testFile(t, "test-file-2-", "test-data")
	defer deleteTestFile(f1)
	defer deleteTestFile(f2)

	v1, v2, err := maybeFromFiles("file://"+f1.Name(), "file://"+f2.Name())
	assert.Nil(t, err)
	assert.Equal(t, "test-pw", v1)
	assert.Equal(t, "test-data", v2)
}

func Test_maybeFromFile(t *testing.T) {
	f := testFile(t, "test-file-x-", "test")
	defer deleteTestFile(f)

	prefixes := []string{"file://", "file://localhost"}
	for _, prefix := range prefixes {
		v, err := maybeFromFile(prefix + f.Name())
		assert.Nil(t, err)
		assert.Equal(t, "test", v)
	}
}

func Test_maybeFromFileErr(t *testing.T) {
	filePath := "file:///var/test/not-existing.txt"
	_, err := maybeFromFile(filePath)
	assert.NotNil(t, err)
}

func Test_maybeFromFileDirect(t *testing.T) {
	s := "test"
	v, err := maybeFromFile(s)
	assert.Nil(t, err)
	assert.Equal(t, "test", v)
}

func testFile(t *testing.T, name string, content string) *os.File {
	f, err := ioutil.TempFile("", name)
	assert.Nil(t, err)
	_, err = f.WriteString(content)
	assert.Nil(t, err)
	return f
}

func deleteTestFile(file *os.File) {
	err := os.Remove(file.Name())
	if err != nil {
		panic(err)
	}
}

func Test_absoluteFilePaths(t *testing.T) {
	v, err := absoluteFilePath("file:///var/test/file.txt")
	assert.Nil(t, err, "absolute path")
	assert.Equal(t, "/var/test/file.txt", v)

	v, err = absoluteFilePath("file://./var/test/file.txt")
	assert.Nil(t, err, "relative path")
	assert.Equal(t, "./var/test/file.txt", v)

	v, err = absoluteFilePath("file://~/var/test/file.txt")
	assert.Nil(t, err, "user directory path")
	assert.True(t, strings.HasSuffix(v, "/var/test/file.txt"))
	assert.False(t, strings.HasSuffix(v, "~/var/test/file.txt"))
	userHomeDir, _ := os.UserHomeDir()
	assert.Equal(t, userHomeDir+"/var/test/file.txt", v)
}

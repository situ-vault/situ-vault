package vault

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maybeFromFiles(t *testing.T) {
	f1 := testFile(t, "test-file-1-", "test-pw")
	f2 := testFile(t, "test-file-2-", "test-data")

	v1, v2, err := maybeFromFiles("file://"+f1.Name(), "file://"+f2.Name())
	assert.Nil(t, err)
	assert.Equal(t, "test-pw", v1)
	assert.Equal(t, "test-data", v2)
}

func Test_maybeFromFile(t *testing.T) {
	f := testFile(t, "test-file-x-", "test")

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

package jsonfile

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultResourceFile(t *testing.T) {
	root := t.TempDir()

	w, err := NewWriter(root)
	assert.Nil(t, err)

	err = w.Write(struct {
		A string `json:"a"`
		B int    `json:"b"`
	}{"hello", 42})

	assert.Equal(t,
		`{
	"a": "hello",
	"b": 42
}`,
		readFile(t, filepath.Join(root, "resource.json")))
}

func TestCustomResourceFile(t *testing.T) {
	root := t.TempDir()

	w, err := NewWriter(root, WithOutputFile("my.json"))
	assert.Nil(t, err)

	err = w.Write(struct {
		Field        string `json:"field"`
		AnotherFiled int    `json:"another-field"`
	}{"hello", 42})

	assert.Equal(t,
		`{
	"field": "hello",
	"another-field": 42
}`,
		readFile(t, filepath.Join(root, "my.json")))
}

func TestCustomIndent(t *testing.T) {
	root := t.TempDir()

	w, err := NewWriter(root, WithIndentation(" "), WithOutputFile("my.json"))
	assert.Nil(t, err)

	err = w.Write(struct {
		Field        string `json:"field"`
		AnotherFiled int    `json:"another-field"`
	}{"hello", 42})

	assert.Equal(t,
		`{
 "field": "hello",
 "another-field": 42
}`,
		readFile(t, filepath.Join(root, "my.json")))
}

func TestConstructorErrors(t *testing.T) {
	root := t.TempDir()

	_, err := NewWriter(root, WithOutputFile("/my.json"))
	assert.ErrorContains(t,
		err,
		"unable to create the Json File Writer: the file is not a relative filename /my.json")

	_, err = NewWriter(root, WithOutputFile("../my.json"))
	assert.ErrorContains(t,
		err,
		"unable to create the Json File Writer: the file is not a relative filename ../my.json")

	_, err = NewWriter(root, WithOutputFile("./my.json"))
	assert.Nil(t, err)
}

func readFile(t *testing.T, f string) string {
	b, err := os.ReadFile(f)
	assert.Nil(t, err, "unable to read file %s", f)
	return string(b)
}

package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCASPathTransformerFunc(t *testing.T) {
	key := "mombestpictures"
	expectedPath := "843b2/03d63/30073/4de10/84eaf/9fae2/7d436/b6c4d"
	expectedOrin := "843b203d63300734de1084eaf9fae27d436b6c4d"

	path := CASPathTransformFunc(key)

	assert.Equal(t, expectedPath, path.PathName)
	assert.Equal(t, expectedOrin, path.Filename)
}
func TestStore(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	key := "momsspecials"
	expectedContent := []byte("some pictures")
	store := NewStorage(opts)
	reader := bytes.NewReader(expectedContent)

	assert.Nil(t, store.writeStream(key, reader))
	r, err := store.Read(key)
	assert.Nil(t, err)
	content, err := io.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	key := "momsspecials"
	content := []byte("some pictures")
	store := NewStorage(opts)
	assert.Nil(t, store.writeStream(key, bytes.NewReader(content)))
	assert.True(t, store.Has(key))
	assert.Nil(t, store.Delete(key))
	assert.False(t, store.Has(key))

}

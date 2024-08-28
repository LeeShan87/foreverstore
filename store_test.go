package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCASPathTransformerFunc(t *testing.T) {
	key := "mombestpictures"
	expectedPath := "843b2/03d63/30073/4de10/84eaf/9fae2/7d436/b6c4d"
	expectedOrin := "843b203d63300734de1084eaf9fae27d436b6c4d"

	path := CASPathTransformFunc(key)

	assert.Equal(t, expectedPath, path.PathName)
	assert.Equal(t, expectedOrin, path.Original)
}
func TestStore(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	store := NewStorage(opts)
	reader := bytes.NewReader([]byte("some pictures"))
	assert.Nil(t, store.writeStream("test_tmp", reader))

}

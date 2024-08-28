package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

type PathTransformFunc func(string) PathKey
type PathKey struct {
	PathName string
	Original string
}

func (p *PathKey) FileName() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Original)
}

type StoreOps struct {
	PathTransformFunc PathTransformFunc
}

type Storage struct {
	StoreOps
}

func NewStorage(opts StoreOps) *Storage {
	return &Storage{
		StoreOps: opts,
	}
}

var DefaultPathTransformFunc = func(key string) PathKey {
	// Really?...
	return PathKey{
		PathName: key,
		Original: key,
	}
}

func (s *Storage) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := pathKey.FileName()
	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("written (%d) bytes to disc: %s\n", n, pathAndFilename)
	return nil
}

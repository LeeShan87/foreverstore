package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey
type PathKey struct {
	PathName string
	Filename string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
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
		Filename: key,
	}
}

func (s *Storage) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)

	fi, err := os.Stat(pathKey.FullPath())
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}
func (s *Storage) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	paths := strings.Split(pathKey.FullPath(), "/")
	for len(paths) > 0 {
		rmPath := strings.Join(paths, "/")
		if err := os.Remove(rmPath); err != nil {
			return err
		}
		//fmt.Println(rmPath)
		paths = paths[:len(paths)-1]
	}
	return nil
}

func (s *Storage) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Storage) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(pathKey.FullPath())
}

func (s *Storage) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := pathKey.FullPath()
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

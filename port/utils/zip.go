package utils

import (
	"bytes"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/klauspost/compress/zip"
)

var blacklist = []string{".DS_Store", "desktop.ini", "Thumbs.db"}

func Zip(in map[string][]byte) ([]byte, error) {
	writer := bytes.NewBuffer([]byte{})
	w := zip.NewWriter(writer)
	for file, data := range in {
		f, err := w.Create(strings.TrimPrefix(file, "/"))
		if err != nil {
			return []byte{}, err
		}
		_, err = f.Write(data)
		if err != nil {
			return []byte{}, err
		}
	}
	w.Close()
	return writer.Bytes(), nil
}

func Unzip(source []byte) (map[string][]byte, error) {
	var out = make(map[string][]byte)
	read, err := zip.NewReader(bytes.NewReader(source), int64(len(source)))
	if err != nil {
		return map[string][]byte{}, err
	}
	var single = true
	var last = ""
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		if slices.Contains(blacklist, filepath.Base(file.Name)) {
			continue
		}
		if strings.Contains(file.Name, "__MACOSX") {
			continue
		}
		name := strings.TrimPrefix(file.Name, "/")
		base := strings.Split(name, "/")[0]
		if single && last != "" && base != last {
			single = false
		}
		last = base
		open, err := file.Open()
		if err != nil {
			return map[string][]byte{}, err
		}
		bytes, err := io.ReadAll(open)
		if err != nil {
			return map[string][]byte{}, err
		}
		out[name] = bytes
	}
	if single {
		oout := out
		out = make(map[string][]byte)
		for name, data := range oout {
			newBase := strings.Join(strings.Split(name, "/")[1:], "/")
			out[newBase] = data
		}
	}
	return out, nil
}

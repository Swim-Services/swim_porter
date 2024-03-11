package utils

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

type PackMeta struct {
	Pack struct {
		Format      int    `json:"pack_format"`
		Description string `json:"description"`
	}
}

func PackMcmeta(mcmeta []byte) (PackMeta, error) {
	mcmeta = []byte(fix(strings.ReplaceAll(string(bytes.TrimPrefix(mcmeta, []byte("\xef\xbb\xbf"))), "\r", "")))
	meta := PackMeta{}
	err := json.Unmarshal(mcmeta, &meta)
	if err != nil {
		return PackMeta{}, err
	}
	return meta, nil
}

var octalEscapePat = regexp.MustCompile(`\\[0-7]{3}`)
var newlineEscape = regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)

func fix(src string) string {
	src = newlineEscape.ReplaceAllStringFunc(src, func(s string) string {
		return strings.ReplaceAll(s, "\n", "\\n")
	})
	return octalEscapePat.ReplaceAllString(src, " ")
}

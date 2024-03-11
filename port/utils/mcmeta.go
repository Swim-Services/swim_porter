package utils

import (
	"bytes"
	"encoding/json"
)

type PackMeta struct {
	Pack struct {
		Format      int    `json:"pack_format"`
		Description string `json:"description"`
	}
}

func PackMcmeta(mcmeta []byte) (PackMeta, error) {
	mcmeta = bytes.TrimPrefix(mcmeta, []byte("\xef\xbb\xbf"))
	meta := PackMeta{}
	err := json.Unmarshal(mcmeta, &meta)
	if err != nil {
		return PackMeta{}, err
	}
	return meta, nil
}

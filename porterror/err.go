package porterror

import "errors"

var (
	ErrMcmetaNotFound   = errors.New("pack.mcmeta not found")
	ErrManifestNotFound = errors.New("manifest.json not found")
)

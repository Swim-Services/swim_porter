package resource

import (
	"encoding/json"

	gtresource "github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/swim-services/swim_porter/jsonnewline"
	stripjsoncomments "github.com/trapcodeio/go-strip-json-comments"
)

// Manifest contains all the basic information about the pack that Minecraft needs to identify it.
type Manifest struct {
	// FormatVersion defines the current version of the manifest. This is currently always 2.
	FormatVersion int `json:"format_version"`
	// Header is the header of a resource pack. It contains information that applies to the entire resource
	// pack, such as the name of the resource pack.
	Header gtresource.Header `json:"header"`
	// Modules describes the modules that comprise the pack. Each entry here defines one of the kinds of
	// contents of the pack.
	Modules []gtresource.Module `json:"modules"`
	// Dependencies describes the packs that this pack depends on in order to work.
	Dependencies []gtresource.Dependency `json:"dependencies,omitempty"`
	// Capabilities are the different features that the pack makes use of that aren't necessarily enabled by
	// default. For a list of options, see below.
	Capabilities []gtresource.Capability `json:"capabilities,omitempty"`

	// Subpacks are the packs's subpacks.
	Subpacks []Subpack `json:"subpacks,omitempty"`
}

type Subpack struct {
	FolderName string `json:"folder_name,omitempty"`
	Name       string `json:"name,omitempty"`
	MemoryTier int    `json:"memory_tier,omitempty"`
}

func UnmarshalJSON(jsonStr []byte) (Manifest, error) {
	commentStripped := stripjsoncomments.Strip(string(jsonStr))
	newlineFixed := jsonnewline.NewLineToEscape(commentStripped)
	manifest := Manifest{}
	if err := json.Unmarshal([]byte(newlineFixed), &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

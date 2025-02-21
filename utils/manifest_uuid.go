package utils

import (
	"github.com/google/uuid"
	"github.com/swim-services/swim_porter/resource"
)

func ChangeUUID(manifest *resource.Manifest) {
	manifest.Header.UUID = uuid.New()
	for i := range manifest.Modules {
		manifest.Modules[i].UUID = uuid.NewString()
	}
}

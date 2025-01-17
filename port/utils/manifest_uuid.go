package utils

import (
	"github.com/google/uuid"
	"github.com/swim-services/swim_porter/port/resource"
)

func ChangeUUID(manifest *resource.Manifest) {
	manifest.Header.UUID = uuid.New()
	for _, module := range manifest.Modules {
		module.UUID = uuid.NewString()
	}
}

package utils

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

func ChangeUUID(manifest *resource.Manifest) {
	manifest.Header.UUID = uuid.NewString()
	for _, module := range manifest.Modules {
		module.UUID = uuid.NewString()
	}
}

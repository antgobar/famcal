package famcal

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed static/*
var embedFrontend embed.FS

func GetFrontendAssets() fs.FS {
	files, err := fs.Sub(embedFrontend, "static")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return files
}

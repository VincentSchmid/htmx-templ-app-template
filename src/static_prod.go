//go:build !dev

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed public/*
var embeddedFiles embed.FS

func publicHandler() http.Handler {
	fsys, err := fs.Sub(embeddedFiles, "public")
	if err != nil {
		panic("failed to locate embedded files")
	}
	return http.StripPrefix("/public/", http.FileServer(http.FS(fsys)))
}

package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static static/* static/vendor/* static/fonts/*
var embeddedStaticFS embed.FS

func embeddedStaticHandler() http.Handler {
	sub, err := fs.Sub(embeddedStaticFS, "static")
	if err != nil {
		return http.NotFoundHandler()
	}
	return http.FileServer(http.FS(sub))
}

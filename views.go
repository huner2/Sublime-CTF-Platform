package main

import (
	"net/http"

	"github.com/flosch/pongo2"
)

var frame = pongo2.Must(pongo2.FromFile("./templates/frame.html")) // Only frame can be pre-compiled from what I can tell

func indexView(config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Explicitly set content-type
	ctx := pongo2.Context{
		"title": config.ctfPrefs.title,
		"page":  "index.html",
	}
	if err := frame.ExecuteWriter(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

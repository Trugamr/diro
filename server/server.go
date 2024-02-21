package server

import (
	"net/http"
)

type server struct {
	root    http.FileSystem
	handler http.Handler
	single  bool
}

func New(root http.FileSystem, single bool) *server {
	return &server{root: root, handler: http.FileServer(root), single: single}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.single {
		// Single mode can only be used if root has an index.html file
		if index, err := s.root.Open("/index.html"); err == nil {
			index.Close()

			// If the requested path is not found or is a directory we rewrite it to root
			f, err := s.root.Open(r.URL.Path)
			if err != nil {
				r.URL.Path = "/"
			} else {
				if stats, err := f.Stat(); err == nil && stats.IsDir() {
					r.URL.Path = "/"
				}
			}
		}
	}

	s.handler.ServeHTTP(w, r)
}

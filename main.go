package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Usage: diro <path>")
	}

	path := os.Args[1]

	// Check if path exists and is a directory
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Path %s does not exist", path)
		}
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatalf("%s is not a directory", path)
	}

	addr := ":8080"

	// Try to listen on the address
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a file server
	server := http.FileServer(http.Dir(path))
	http.Handle("/", server)

	log.Printf("Server started on http://localhost%s", addr)

	// Serve the files
	log.Fatal(http.Serve(l, nil))
}

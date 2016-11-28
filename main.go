package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Check that we're given a serial device name
	if len(os.Args) < 2 {
		fmt.Println("expected a serial device name")
	}
	serialName := os.Args[1]

	// Connect to the device
	d, err := newDevice(serialName)
	if err != nil {
		log.Fatal(err)
	}

	// Set up a static directory for resources
	fs := http.FileServer(http.Dir("resources"))
	http.Handle("/resources/", http.StripPrefix("/resources", fs))

	// Handle data requests
	http.Handle("/data", &dataHandler{d})

	// Handle the homepage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pages/home.html")
	})

	// Start the server
	log.Fatalln(http.ListenAndServe(":8080", nil))

}

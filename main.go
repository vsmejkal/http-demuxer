package main

import (
	"log"
	"os"
	"net/http"
	"fmt"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		log.Fatalf("Usage: %s routes.json\n", os.Args[0])
	}

	// Load routes from file
	var conf Config
	if err := conf.Load(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	// Start server
	addr := fmt.Sprintf(":%d", conf.Port)
	log.Println("Listening on port", conf.Port, "...")
	log.Fatal(http.ListenAndServe(addr, NewDemuxer(conf)))
}

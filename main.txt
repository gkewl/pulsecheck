package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello I'm Automatically deployed using Travis")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

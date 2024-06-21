//go:build debug

package main

import (
	"fmt"
	"log"
	"net/http"
)

// init initializes the pprof server.
func init() {
	pprof = func() {
		go func() {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "")
			})
			log.Fatal(http.ListenAndServe(":80", nil))
		}()
	}
}

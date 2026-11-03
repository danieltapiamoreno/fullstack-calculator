// Command server runs the calculator API.
package main

import (
	"log"
	"net/http"

	"fullstack-calculator/backend/internal/api"
)

func main() {
	const address = ":8080"
	log.Printf("listening on %s", address)
	log.Fatal(http.ListenAndServe(address, api.NewHandler()))
}

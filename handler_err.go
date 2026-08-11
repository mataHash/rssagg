package main

import (
	"net/http"
)

func handler_err(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}

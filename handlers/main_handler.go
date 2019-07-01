package handlers

import (
	"fmt"
	"net/http"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Supported Routes: %s", r.URL.Path[1:])
}

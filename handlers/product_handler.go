package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func HandleProduct(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Fprintf(w, "Supported Routes: %s", r.URL.Path[1:])
}

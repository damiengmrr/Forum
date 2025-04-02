package handlers

import (
	"fmt"
	"net/http"
)

func TestSessionHandler(w http.ResponseWriter, r *http.Request) {
	id := GetCurrentUserID(r)
	fmt.Fprintf(w, "Ton ID : %d", id)
}

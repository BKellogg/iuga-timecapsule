package handlers

import (
	"net/http"
)

// HandleNewTimeCapsule handles a new submission of a time capsule
func HandleNewTimeCapsule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "all requests to create a new time capsule must be POST", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "requests must have a Content-Type of application/json", http.StatusBadRequest)
		return
	}
}

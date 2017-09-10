package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BKellogg/iuga-timecapsule/tc-server/models"
)

// HandleNewTimeCapsule handles a new submission of a time capsule
func (ctx *HandlerContext) HandleNewTimeCapsule(w http.ResponseWriter, r *http.Request) {
	// ensure a proper request
	if r.Method != "POST" {
		http.Error(w, "all requests to create a new time capsule must be POST", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "requests must have a Content-Type of application/json", http.StatusBadRequest)
		return
	}

	// decode the request body inot a NewCapsule
	capsule := &models.NewCapsule{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(capsule); err != nil {
		http.Error(w, "error decoding request into a new capsule"+err.Error(), http.StatusInternalServerError)
		return
	}

	// insert the new capsule
	db := ctx.DB
	_, err := db.Exec(`INSERT INTO Capsules (NetID, GradDate, Message) VALUES (` +
		`'` + capsule.NetID + `', '` + capsule.GradDate + `', '` + capsule.Message + `')`)
	if err != nil {
		http.Error(w, "error inserting new capsule: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("success!"))
}

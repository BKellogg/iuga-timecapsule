package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BKellogg/iuga-timecapsule/tc-server/models"
	"github.com/BKellogg/iuga-timecapsule/tc-server/utils"
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
	if len(capsule.NetID) == 0 || len(capsule.GradDate) == 0 || len(capsule.Message) == 0 {
		http.Error(w, "all inputs are required", http.StatusBadRequest)
		return
	}

	// check if the user has already submitted a capsule and stop the request if the user has submitted one
	if err := userHasAlreadySubmittedCapsule(ctx.DB, capsule.NetID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// insert the new capsule
	err := insertNewCapsule(ctx.DB, capsule)
	if err != nil {
		http.Error(w, "error inserting new capsule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send a success email with the message the user sent in
	if err = utils.SendSuccessEmail(ctx.GmailService, capsule); err != nil {
		http.Error(w, "error sending confirmation email"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	response := models.CapsuleResponse{
		GradDate: capsule.GradDate,
		Email:    capsule.NetID + "@uw.edu",
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

// returns if the user has already submitted a capsule, and an error if one occured
func userHasAlreadySubmittedCapsule(db *sql.DB, netid string) error {
	result, err := db.Query(`SELECT * FROM Capsules WHERE Capsules.NetID=?`, netid)
	if err != nil {
		return err
	}
	if result.Next() {
		return errors.New("time capsule already found for this user")
	}
	return nil
}

// inserts the given new capsule into the given db, returns an error if there was one
func insertNewCapsule(db *sql.DB, capsule *models.NewCapsule) error {
	_, err := db.Exec(`INSERT INTO Capsules (NetID, GradDate, Message) VALUES (` +
		`'` + capsule.NetID + `', '` + capsule.GradDate + `', '` + capsule.Message + `')`)
	return err
}

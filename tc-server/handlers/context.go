package handlers

import (
	"database/sql"

	gmail "google.golang.org/api/gmail/v1"
)

// HandlerContext defines context for handlers that includes
// any global variables they need.
type HandlerContext struct {
	DB           *sql.DB
	GmailService *gmail.Service
}

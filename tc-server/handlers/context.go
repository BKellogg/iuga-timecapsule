package handlers

import "database/sql"

// HandlerContext defines context for handlers that includes
// any global variables they need.
type HandlerContext struct {
	DB *sql.DB
}

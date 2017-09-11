package handlers

import "net/http"

// CORSHandler defines a CORS wrapper for a http.Handler
type CORSHandler struct {
	Handler http.Handler
}

func (ch *CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization")

	//if this is preflight request, the method will
	//be OPTIONS, so call the real handler only if
	//the method is something else
	if r.Method != "OPTIONS" {
		ch.Handler.ServeHTTP(w, r)
	}
}

// NewCORSHandler wrappers the given http.Handler in CORS
func NewCORSHandler(handlerToWrap http.Handler) *CORSHandler {
	return &CORSHandler{handlerToWrap}
}

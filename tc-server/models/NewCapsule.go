package models

// NewCapsule represents the format of a request
// made to create a new Capsule
type NewCapsule struct {
	NetID    string `json:"netID"`
	GradDate string `json:"gradDate"`
	Message  string `json:"message"`
}

package dto

// ErrorResponse is the standard JSON error response body returned by all
// API endpoints on failure.
type ErrorResponse struct {
	Error string `json:"error"`
}

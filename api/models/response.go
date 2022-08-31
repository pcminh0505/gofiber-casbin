package models

// Response model for Swagger API
type Response struct {
	Error   bool
	Message string
	Data    interface{} `json:",omitempty"`
}

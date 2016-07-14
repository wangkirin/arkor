package models

type Error struct {
	Type      string `json:"type,omitempty"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Resource  string `json:"resource,omitempty"`
	RequestID string `json:"requestID,omitempty"`
}

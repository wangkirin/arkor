package models

type Error struct {
	Type      string `json:"Type,omitempty"`
	Code      string `json:"Code,omitempty"`
	Message   string `json:"Message,omitempty"`
	Resource  string `json:"Resource,omitempty"`
	RequestID string `json:"RequestId,omitempty"`
}

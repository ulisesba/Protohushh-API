// Package protohush provides structures and interfaces related to chat intentions.
package protohush

// ChatQuery represents a structure to define user's intention during a chat.
type ChatQuery struct {
	// Intention specifies the main purpose or action the user intends to perform.
	Intention string

	// Value represents the specific term or keyword that the user intends to search for.
	// The json tag indicates that during JSON encoding or decoding, this field maps to "value_to_search".
	Value string `json:"value_to_search"`

	// Collection indicates the main set or group where the search should be made.
	Collection string

	Limit int `json:"limit"`
}

type ChatResponse struct {
	Intention string `json:"intention"`
	Data      string `json:"data"`
}

// Chat defines an interface for chat-related operations.
type Chat interface {
	// Search is a method that takes a value as an argument and performs a search operation based on it.
	Search(instruction string) (*ChatResponse, error)
}

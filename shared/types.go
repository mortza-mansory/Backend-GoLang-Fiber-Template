package shared

import "time"

// ID is the primary identifier used across the domain. Using a named string
// type makes signatures clearer without introducing a framework.
type ID string

// IDGenerator produces unique identifiers (e.g. UUIDs).
type IDGenerator interface {
	NewID() ID
}

// Auditable captures common created/updated timestamps. Modules may embed it.
type Auditable struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

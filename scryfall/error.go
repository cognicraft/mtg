package scryfall

// An Error object represents a failure to find information or understand the input you provided to the API.
// Error objects are always transmitted with the appropriate 4XX or 5XX HTTP status code.
type Error struct {
	Object string `json:"object"`

	// An integer HTTP status code for this error.
	Status int `json:"status"`

	// A computer-friendly string representing the appropriate HTTP status code.
	Code string `json:"code"`

	// A human-readable string explaining the error.
	Details string `json:"details"`

	// A computer-friendly string that provides additional context for the main error.
	// For example, an endpoint many generate HTTP 404 errors for different kinds of input.
	// This field will provide a label for the specific kind of 404 failure, such as ambiguous.
	Type string `json:"type,omitempty"`

	// If your input also generated non-failure warnings, they will be provided as human-readable strings in this array.
	Warnings []string `json:"warnings"`
}

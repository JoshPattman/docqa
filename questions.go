package docqa

// Question defines a sepcific question to send to the LLM.
type Question struct {
	// Question is the one-liner, for example `Who is the author of this document?`.
	Question string `json:"question"`
	// Details provides extra details, examples, and instructions on how to answer this question, as bullet points.
	Details []string `json:"details"`
	// AllowedTypeKeys lists all the keys of the types which the LLM may respond with.
	AllowedTypeKeys []string `json:"allowed_type_keys"`
}

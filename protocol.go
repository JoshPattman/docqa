package docqa

// Protocol defines a method of communication to and from the LLM.
type Protocol interface {
	// Schema creates a jsonschema given a set of keyed questions.
	Schema(qs map[string]Question) map[string]any
	// SystemPrompt builds a system prompt for the given set of [Question]s.
	SystemPrompt(qs map[string]Question) string
	// ParseResponse takes a raw LLM response and parses it into lists of [Entity], keyed by question key.
	ParseResponse(resp string) (map[string][]Entity, error)
}

// RoleAndTask defines a role and a task for the LLM,
// both of which should be short (~ 1 sentence).
type RoleAndTask struct {
	Role string
	Task string
}

// ExtractAnswers answers the given [Question]s about a document,
// returning lists of [Entity] keyed by question key.
func ExtractAnswers(client Client, qa Protocol, questions map[string]Question, documentText string) (map[string][]Entity, LLMUsage, error) {
	resp, usage, err := client.GetLLMResponse(
		qa.SystemPrompt(questions),
		documentText,
		qa.Schema(questions),
	)
	if err != nil {
		return nil, usage, err
	}

	answers, err := qa.ParseResponse(resp)
	if err != nil {
		return nil, usage, err
	}
	return answers, usage, nil
}

// GetDefaultRoleAndTask builds a [RoleAndTask] for a generic document information extraction task.
func GetDefaultRoleAndTask() RoleAndTask {
	return RoleAndTask{
		Role: "You are a state-of-the-art data extraction AI",
		Task: "Your task is to answer the provided questions completely accurately, based on the user-provided document",
	}
}

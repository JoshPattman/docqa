package docqa

// Wrap a standard json schema into one that can be sent to openai
func wrapOpenAISchema(schema map[string]any) map[string]any {
	return map[string]any{
		"type": "json_schema",
		"json_schema": map[string]any{
			"name":   "qa-schema",
			"strict": true,
			"schema": schema,
		},
	}
}

package docqa

import (
	"encoding/json"
	"fmt"
	"strings"
)

type basicProtocol struct {
	roleAndTask RoleAndTask
	types       map[string]Type
}

// NewBasicProtocol creates a sensible and generalised protocol for information extraction.
func NewBasicProtocol(roleTask RoleAndTask, types map[string]Type) Protocol {
	return &basicProtocol{
		roleAndTask: roleTask,
		types:       types,
	}
}

// Schema implements [Protocol].
func (qa *basicProtocol) Schema(qs map[string]Question) map[string]any {
	properties := make(map[string]any)
	components := make(map[string]any)

	for key, parser := range qa.types {
		schemaProps := parser.SchemaProperties()
		schemaProps["answer_type"] = map[string]any{
			"const": key,
			"type":  "string",
		}
		schemaPropsKeys := make([]string, 0, len(schemaProps))
		for k := range schemaProps {
			schemaPropsKeys = append(schemaPropsKeys, k)
		}
		components[key] = map[string]any{
			"type":                 "object",
			"properties":           schemaProps,
			"required":             schemaPropsKeys,
			"additionalProperties": false,
		}
	}

	for key, question := range qs {
		options := make([]map[string]any, 0)
		for _, tk := range question.AllowedTypeKeys {
			options = append(options, map[string]any{
				"$ref": fmt.Sprintf("#/definitions/%s", tk),
			})
		}
		one := map[string]any{
			"anyOf": options,
		}
		properties[key] = map[string]any{
			"type":  "array",
			"items": one,
		}
	}

	pkeys := make([]string, 0, len(properties))
	for k := range properties {
		pkeys = append(pkeys, k)
	}

	return map[string]any{
		"type":                 "object",
		"properties":           properties,
		"definitions":          components,
		"required":             pkeys,
		"additionalProperties": false,
	}
}

// ParseResponse implements [Protocol].
func (qa *basicProtocol) ParseResponse(resp string) (map[string][]Entity, error) {
	respTyped := make(map[string][]map[string]any)
	err := json.Unmarshal([]byte(resp), &respTyped)
	if err != nil {
		return nil, err
	}
	answers := make(map[string][]Entity)
	for qKey, qAnswers := range respTyped {
		answers[qKey] = []Entity{}
		for _, qAnswer := range qAnswers {
			answerType, ok := qAnswer["answer_type"]
			if !ok {
				return nil, fmt.Errorf("answer did not have answer_type key")
			}
			answerTypeStr, ok := answerType.(string)
			if !ok {
				return nil, fmt.Errorf("answer type was not a string")
			}
			parser, ok := qa.types[answerTypeStr]
			if !ok {
				return nil, fmt.Errorf("did not have a parser for that")
			}
			entity, err := parser.Parse(qAnswer)
			if err != nil {
				return nil, err
			}
			entity.Attr().LocalisedRange = IndefRange()
			entity.Attr().EvidenceRanges = make([]Range, 0)
			answers[qKey] = append(answers[qKey], entity)
		}
	}
	return answers, nil
}

// SystemPrompt implements [Protocol].
func (qa *basicProtocol) SystemPrompt(qs map[string]Question) string {
	builder := &mdBuilder{}
	builder.Headerf(1, "Role & Task")
	builder.Bullet(0, qa.roleAndTask.Role)
	builder.Bullet(0, qa.roleAndTask.Task)
	builder.Bulletf(0, "The user will provide you with the raw text from the document in question")
	builder.Break(2)

	builder.Headerf(1, "Answer Types")
	builder.Bulletf(0, "You can answer each question with some amount of answer objects")
	builder.Bulletf(0, "Each type of answer object has a different purpose, with different properties")
	builder.Bulletf(0, "Below are the allowed answer types")
	for key, parser := range qa.types {
		builder.Break(1)
		builder.Headerf(2, "`%s`", key)
		instructions := parser.Instructions()
		builder.Bulletf(0, "**%s**", instructions.OneLiner)
		for _, d := range instructions.Details {
			builder.Bullet(0, d)
		}
	}
	builder.Break(2)

	builder.Headerf(1, "Questions")
	builder.Bulletf(0, "You should answer all questions")
	builder.Bulletf(0, "If you cannot determine the answer to a question, return an empty list for that question")
	for key, question := range qs {
		builder.Break(1)
		builder.Headerf(1, "`%s`", key)
		builder.Bulletf(0, "**%s**", question.Question)
		for _, d := range question.Details {
			builder.Bullet(0, d)
		}
		builder.Bulletf(0, "Allowed response types: %s", strings.Join(question.AllowedTypeKeys, ", "))
	}

	return builder.Build()
}

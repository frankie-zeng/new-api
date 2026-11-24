package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionParametersUsesJSONSchemaForVertexAnyOf(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"moduleName": map[string]any{
				"description": "module",
				"anyOf": []any{
					map[string]any{"type": "string"},
					map[string]any{"type": "object", "properties": map[string]any{
						"id": map[string]any{"type": "string"},
					}},
				},
			},
		},
	}

	field, value, misfits := FunctionParameters(schema, true)
	assert.Equal(t, "parametersJsonSchema", field)
	assert.Empty(t, misfits)
	assert.Equal(t, schema, value)
}

func TestFunctionParametersUsesJSONSchemaForMissingType(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"values": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"value": map[string]any{"description": "any value"},
				},
			},
		},
	}

	field, value, misfits := FunctionParameters(schema, true)
	assert.Equal(t, "parametersJsonSchema", field)
	assert.Empty(t, misfits)
	assert.Equal(t, schema, value)
}

func TestFunctionParametersKeepsATypedSchema(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	field, _, misfits := FunctionParameters(schema, true)
	assert.Equal(t, "parameters", field)
	assert.Empty(t, misfits)
}

func TestFunctionParametersLegacyCollapsesNullableAnyOf(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"surface": map[string]any{
				"description": "target surface",
				"anyOf": []any{
					map[string]any{"type": "string", "enum": []any{"word", "excel"}},
					map[string]any{"type": "null"},
				},
			},
		},
	}

	field, value, misfits := FunctionParameters(schema, false)
	assert.Equal(t, "parameters", field)
	assert.Contains(t, misfits, "null type")

	cleaned, ok := value.(map[string]any)
	require.True(t, ok)
	surface := cleaned["properties"].(map[string]any)["surface"].(map[string]any)
	assert.Equal(t, "STRING", surface["type"])
	assert.Equal(t, true, surface["nullable"])
	assert.Equal(t, "target surface", surface["description"])
	assert.Equal(t, []any{"word", "excel"}, surface["enum"])
	assert.NotContains(t, surface, "anyOf")
}

package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanFunctionParametersFillsMissingSchemaType(t *testing.T) {
	cleaned, ok := CleanFunctionParameters(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"values": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"value": map[string]any{"description": "any value"},
					"empty": map[string]any{},
					"none":  map[string]any{"type": "null"},
					"list":  map[string]any{"items": map[string]any{"type": "string"}},
					"bag":   map[string]any{"properties": map[string]any{"id": map[string]any{"type": "string"}}},
				},
			},
		},
	}).(map[string]any)
	require.True(t, ok)

	values := cleaned["properties"].(map[string]any)["values"].(map[string]any)
	fields := values["properties"].(map[string]any)

	assert.Equal(t, "STRING", fields["value"].(map[string]any)["type"])
	assert.Equal(t, "STRING", fields["empty"].(map[string]any)["type"])

	none := fields["none"].(map[string]any)
	assert.Equal(t, "STRING", none["type"])
	assert.Equal(t, true, none["nullable"])

	list := fields["list"].(map[string]any)
	assert.Equal(t, "ARRAY", list["type"])
	assert.Equal(t, "STRING", list["items"].(map[string]any)["type"])

	bag := fields["bag"].(map[string]any)
	assert.Equal(t, "OBJECT", bag["type"])
	assert.Equal(t, "STRING", bag["properties"].(map[string]any)["id"].(map[string]any)["type"])
}

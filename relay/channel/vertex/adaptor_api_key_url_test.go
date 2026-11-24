package vertex

import (
	"net/http/httptest"
	"testing"

	"github.com/MoeForge/moego/ojson"
	"github.com/QuantumNous/new-api/common"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRequestURLWithAPIKey(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		region    string
		stream    bool
		want      string
	}{
		{
			name:      "project ID keeps the global region",
			projectID: "my-project",
			region:    `{"default":"global"}`,
			want:      "https://aiplatform.googleapis.com/v1/projects/my-project/locations/global/publishers/google/models/gemini-4-argon:generateContent?key=test-key",
		},
		{
			name:      "project ID keeps a regional endpoint when streaming",
			projectID: "my-project",
			region:    `{"default":"us-central1"}`,
			stream:    true,
			want:      "https://us-central1-aiplatform.googleapis.com/v1/projects/my-project/locations/us-central1/publishers/google/models/gemini-4-argon:streamGenerateContent?alt=sse&key=test-key",
		},
		{
			name:   "no project ID keeps the express-mode URL",
			region: `{"default":"global"}`,
			want:   "https://aiplatform.googleapis.com/v1/publishers/google/models/gemini-4-argon:generateContent?key=test-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &relaycommon.RelayInfo{
				IsStream:        tt.stream,
				OriginModelName: "gemini-4-argon",
				ChannelMeta: &relaycommon.ChannelMeta{
					ApiVersion:        tt.region,
					ApiKey:            "test-key",
					UpstreamModelName: "gemini-4-argon",
					ChannelOtherSettings: dto.ChannelOtherSettings{
						VertexKeyType:   dto.VertexKeyTypeAPIKey,
						VertexProjectID: tt.projectID,
					},
				},
			}
			adaptor := &Adaptor{}
			adaptor.Init(info)

			got, err := adaptor.GetRequestURL(info)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertClaudeRequestUsesGeminiPayloadForGeminiModels(t *testing.T) {
	body := `{
		"model": "gemini-4-neon",
		"max_tokens": 1024,
		"stream": true,
		"system": "be brief",
		"thinking": {"type": "enabled", "budget_tokens": 1024},
		"messages": [{"role": "user", "content": "hi"}],
		"tools": [{
			"name": "lookup",
			"description": "look up",
			"eager_input_streaming": true,
			"cache_control": {"type": "ephemeral"},
			"input_schema": {"type": "object", "properties": {"q": {"type": "string"}}}
		}]
	}`
	request := &dto.ClaudeRequest{}
	require.NoError(t, common.UnmarshalJsonStr(body, request))

	info := &relaycommon.RelayInfo{
		RelayFormat:     types.RelayFormatClaude,
		OriginModelName: "gemini-4-neon",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "gemini-4-neon",
		},
	}
	adaptor := &Adaptor{}
	adaptor.Init(info)

	converted, err := adaptor.ConvertClaudeRequest(nil, info, request)
	require.NoError(t, err)
	sent, err := common.Marshal(converted)
	require.NoError(t, err)
	wire := ojson.ParseBytes(sent)

	raw := string(sent)
	assert.NotContains(t, raw, "anthropic_version")
	assert.NotContains(t, raw, "input_schema")
	assert.NotContains(t, raw, "eager_input_streaming")
	assert.NotContains(t, raw, "cache_control")
	assert.False(t, wire.Get("messages").Exists())
	assert.False(t, wire.Get("thinking").Exists())
	assert.Equal(t, "user", wire.Get("contents.0.role").String())
	assert.Equal(t, "hi", wire.Get("contents.0.parts.0.text").String())
	assert.Equal(t, "lookup", wire.Get("tools.0.functionDeclarations.0.name").String())
	assert.Equal(t, "OBJECT", wire.Get("tools.0.functionDeclarations.0.parameters.type").String())
	assert.Equal(t, int64(1024), wire.Get("generationConfig.maxOutputTokens").Int())
}

func TestConvertClaudeRequestKeepsAnthropicPayloadForClaudeModels(t *testing.T) {
	request := &dto.ClaudeRequest{}
	require.NoError(t, common.UnmarshalJsonStr(`{
		"model": "claude-sonnet-4-5-20250929",
		"max_tokens": 64,
		"messages": [{"role": "user", "content": "hi"}]
	}`, request))

	info := &relaycommon.RelayInfo{
		RelayFormat:     types.RelayFormatClaude,
		OriginModelName: "claude-sonnet-4-5-20250929",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "claude-sonnet-4-5-20250929",
		},
	}
	adaptor := &Adaptor{}
	adaptor.Init(info)
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	converted, err := adaptor.ConvertClaudeRequest(ctx, info, request)
	require.NoError(t, err)
	sent, err := common.Marshal(converted)
	require.NoError(t, err)
	wire := ojson.ParseBytes(sent)

	assert.Equal(t, "vertex-2023-10-16", wire.Get("anthropic_version").String())
	assert.Equal(t, "hi", wire.Get("messages.0.content").String())
}

package chatgptplugin

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpzouying/chatgptplugin/llm"
	"github.com/xpzouying/chatgptplugin/openai"
	"github.com/xpzouying/chatgptplugin/plugins/calculator"
	"github.com/xpzouying/chatgptplugin/plugins/v2ex"
)

func TestManagerHandle_Calculator(t *testing.T) {

	manager := newChatGPTManager()
	{
		manager.AddPlugin(calculator.NewCalculator())
	}

	answer, err := manager.Handle(context.Background(), "10 add 20 equals ?")
	require.NoError(t, err)

	assert.True(t, answer["result"].(bool))
	want := float64(30)
	assert.Equal(t, want, answer["message"].(float64))
}

func TestManagerHandle_V2ex(t *testing.T) {

	manager := newChatGPTManager()
	{
		manager.AddPlugin(v2ex.NewV2ex())
	}

	answer, err := manager.Handle(context.Background(), "查看 v2ex 热榜列表")
	require.NoError(t, err)

	assert.True(t, answer["result"].(bool))
	assert.NotEmpty(t, answer["data"].(v2ex.HotsList))
}

func newChatGPTManager() *Manager {
	_ = godotenv.Load() // ignore if file not exists

	var llmer llm.LLMer
	{
		token := os.Getenv("OPENAI_TOKEN")
		if len(token) == 0 {
			panic("empty openai token: set os env: OPENAI_TOKEN")
		}
		llmer = openai.NewChatGPT(token, openai.WithModel("gpt-4"))
	}

	return NewManager(llmer)
}

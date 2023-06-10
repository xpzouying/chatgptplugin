package llm

import (
	"context"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

func (r Role) String() string {
	return string(r)
}

type LlmMessage struct {
	Role    Role
	Content string
}

type LlmAnswer struct {
	Role    string
	Content string
}

type LLMer interface {
	Chat(ctx context.Context, messages []LlmMessage) (*LlmAnswer, error)
}

type Summarizer interface {
	Summary(ctx context.Context, content string) (string, error)
}

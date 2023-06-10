package chatgptplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/xpzouying/chatgptplugin/llm"
)

var (
	ErrNoValidPlugin = errors.New("no valid plugin")
)

type PluginContext struct {
	Plugin

	// Request for handle function of plugin.
	Request map[string]any
}

type Manager struct {
	llmer llm.LLMer

	// plugins <key:name, value:Plugin>
	plugins map[string]Plugin
}

type ManagerOpt func(manager *Manager)

// WithPlugin enable one plugin.
func WithPlugin(p Plugin) ManagerOpt {

	return func(manager *Manager) {
		name := strings.ToLower(p.GetName())
		if _, ok := manager.plugins[name]; !ok {
			manager.plugins[name] = p
		}
	}
}

// WithPlugins enable multiple plugins.
func WithPlugins(plugins []Plugin) ManagerOpt {

	return func(manager *Manager) {

		for _, p := range plugins {
			opt := WithPlugin(p)
			opt(manager)
		}
	}
}

// NewManager create plugin manager.
func NewManager(llmer llm.LLMer, opts ...ManagerOpt) *Manager {

	manager := &Manager{
		llmer:   llmer,
		plugins: make(map[string]Plugin, 4),
	}

	for _, opt := range opts {
		opt(manager)
	}

	return manager
}

func (m *Manager) AddPlugin(plugin Plugin) {
	m.plugins[plugin.GetName()] = plugin
}

func (m *Manager) Handle(ctx context.Context, query string) (map[string]any, error) {

	pluginCtx, err := m.Select(ctx, query)
	if err != nil {
		return nil, err
	}

	answer, err := pluginCtx.Do(ctx, pluginCtx.Request)
	if err != nil {
		return nil, err
	}

	log.Printf("got plugin answer: %v", answer)
	return answer, nil
}

// Select to choice some plugin to finish the task.
func (m *Manager) Select(ctx context.Context, query string) (*PluginContext, error) {

	answer, err := m.chatWithLlm(ctx, query)
	if err != nil {
		return nil, err
	}

	return m.choicePlugins(answer)
}

func (m *Manager) makePrompt(query string) string {

	tools := m.makeTaskList()

	prompt := fmt.Sprintf(`你的目标任务是：%s

你有一些插件工具可以选择，如果没有找到合适的插件，则直接返回空的 json 格式 '{}'。
返回调用插件的格式请一定要使用 json 的格式，返回的格式如下：
'''
{
  "plugin": "$PluginName",
  "args": { $ArgsExample }
}
'''
其中，$PluginName 替换成插件的名字，$ArgsExample 替换成插件的参数。
格式里面的 key 参数请保持跟对应示例中的保持一致，不要随意修改 json key 的名字。
当你选择出合适的工具后，请不要解释你为什么选择该工具，只需要告诉我选择的工具以及处理后的参数。

例如：假设用户提供了 Google 的插件以及对应的请求参数示例如下：
'''
* Google: 可以进行网络搜索。请求参数示例为：'{"query": "搜索词"}'
'''

那么，当用户搜索明天是周几时，则应该返回：
'''
{
  "plugin": "Google",
  "args": {
    "query": "明天是周几"
  }
}
'''

如果没有合适的工具，或者你不确定应该选择什么工具完成用户的任务的话，那么返回空的 json 即可，例如：
'''
{}
'''

现在，你可以选择工具有下面的几种工具，根据用户的目标，选择下列的工具中的一个，下面会给出插件的名字、它的作用、以及对应的 json 格式的参数示例：

'''
%s
'''
`,

		query,
		tools,
	)

	return prompt
}

func (m *Manager) makeTaskList() string {

	lines := make([]string, 0, len(m.plugins))

	for _, p := range m.plugins {

		line := fmt.Sprintf(
			`* %s: 该工具的作用是：%s, 请求参数示例为: %s`,
			p.GetName(),
			p.GetInputExample(),
			p.GetDesc(),
		)

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m *Manager) chatWithLlm(ctx context.Context, query string) (string, error) {
	prompt := m.makePrompt(query)

	messages := []llm.LlmMessage{
		{
			Role:    llm.RoleSystem,
			Content: "You are an helpful and kind assistant to answer questions that can use tools to interact with real world and get access to the latest information.",
		},
		{
			Role:    llm.RoleUser,
			Content: prompt,
		},
	}

	answer, err := m.llmer.Chat(ctx, messages)
	if err != nil {
		return "", errors.Wrap(err, "chat with llmer failed")
	}

	return answer.Content, nil
}

func (m *Manager) choicePlugins(answer string) (*PluginContext, error) {

	var pluginAnswer struct {
		Plugin string         `json:"plugin,omitempty"`
		Args   map[string]any `json:"args,omitempty"`
	}

	if err := json.Unmarshal([]byte(answer), &pluginAnswer); err != nil {
		return nil, err
	}

	if pluginAnswer.Plugin == "" {
		return nil, ErrNoValidPlugin
	}

	var (
		name = pluginAnswer.Plugin
		req  = pluginAnswer.Args
	)

	if p, ok := m.plugins[name]; ok {
		pluginCtx := &PluginContext{
			Plugin:  p,
			Request: req,
		}

		return pluginCtx, nil
	}

	return nil, ErrNoValidPlugin
}
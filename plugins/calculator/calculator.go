package calculator

import (
	"context"

	"github.com/mnogu/go-calculator"
	"github.com/pkg/errors"
	"github.com/xpzouying/chatgptplugin/plugins"
)

const (
	pluginName         = `Calculator`
	pluginDesc         = `A calculator, capable of performing mathematical calculations, where the input is a description of a mathematical expression and the return is the result of the calculation. For example: the input is: one plus two, the return is three.`
	pluginInputExample = `{"input": "1+2"}`
)

type Calculator struct{}

func NewCalculator() *Calculator {

	return &Calculator{}
}

func (c Calculator) GetInputExample() string {
	return pluginInputExample
}

func (Calculator) Do(ctx context.Context, req map[string]any) (map[string]any, error) {
	input, ok := req["input"]
	if !ok {
		return nil, plugins.ErrInvalidPluginReq
	}

	s := input.(string)

	result, err := calculator.Calculate(s)
	if err != nil {
		return nil, errors.Wrap(err, "calculate failed")
	}

	return map[string]any{
		"result":  true,
		"message": result,
	}, nil
}

func (c Calculator) GetName() string {
	return pluginName
}

func (c Calculator) GetDesc() string {
	return pluginDesc
}

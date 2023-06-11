package v2ex

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

var (
	httpClient = &http.Client{}
)

const (
	urlV2exHots = "https://www.v2ex.com/api/topics/hot.json"

	pluginName         = "v2ex"
	pluginDesc         = "v2ex 是一个由设计师、程序员及有创意的人参与的社区。"
	pluginInputExample = `{}`
)

type V2ex struct{}

func NewV2ex() *V2ex {
	return &V2ex{}
}

func (v *V2ex) Do(ctx context.Context, req map[string]any) (map[string]any, error) {

	return v.sendRequest(ctx)
}

func (v *V2ex) sendRequest(ctx context.Context) (map[string]any, error) {
	data, err := v.getV2exHotsHttpData(ctx)
	if err != nil {
		return nil, err
	}

	return v.decodeV2exHotsList(data)
}

func (v *V2ex) getV2exHotsHttpData(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, urlV2exHots, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new http request failed")
	}
	req = req.WithContext(ctx)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send http client failed")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read http response failed")
	}

	return data, nil
}

type (
	Hots     map[string]any
	HotsList []Hots
)

func (v *V2ex) decodeV2exHotsList(data []byte) (map[string]any, error) {
	// https://www.v2ex.com/api/topics/hot.json

	var result HotsList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.Wrap(err, "failed decode http response")
	}

	return map[string]any{
		"result": true,
		"data":   result,
	}, nil
}

func (v *V2ex) GetName() string {
	return pluginName
}

func (v *V2ex) GetInputExample() string {
	return pluginInputExample
}

func (v *V2ex) GetDesc() string {
	return pluginDesc
}

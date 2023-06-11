package v2ex

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	v := NewV2ex()

	got, err := v.sendRequest(context.Background())
	assert.NoError(t, err)

	assert.NotEmpty(t, got)
	assert.True(t, got["result"].(bool))

	hots := got["data"].(HotsList)
	assert.NotEmpty(t, hots)

	t.Logf("got v2ex HotsList list: %v", hots)
}

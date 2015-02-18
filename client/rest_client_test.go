package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRestClient(t *testing.T) {
	url := "http://127.0.0.1/test"
	space := "space"
	token := "token"

	client := NewRestClient(url, token, space)

	assert.Equal(t, client.Token, token)
	assert.Equal(t, client.Space, space)
	assert.Equal(t, client.TargetURL, url)

	assert.Panics(t, func() {
		client = NewRestClient(url, "", space)
	})
}

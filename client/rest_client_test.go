package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Successful_NewRestClient(t *testing.T) {
	url := "http://127.0.0.1/test"
	space := "space"
	token := "token"

	client := NewRestClient(url, token, space)

	assert.Equal(t, client.Token, token)
	assert.Equal(t, client.Space, space)
	assert.Equal(t, client.TargetURL, url)
}

func Test_NewRestClient_WithBadToken(t *testing.T) {
	url := "http://127.0.0.1/test"
	space := "space"
	token := ""

	assert.Panics(t, func() {
		_ = NewRestClient(url, token, space)
	})
}

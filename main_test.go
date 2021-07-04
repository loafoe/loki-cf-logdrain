package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListenString(t *testing.T) {
	port := os.Getenv("PORT")
	defer func() {
		_ = os.Setenv("PORT", port)
	}()
	_ = os.Setenv("PORT", "")
	s := listenString()
	assert.Equal(t, s, ":8080")
	_ = os.Setenv("PORT", "1028")
	s = listenString()
	assert.Equal(t, s, ":1028")
}

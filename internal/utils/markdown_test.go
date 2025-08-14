package utils

import (
	"encoding/json"
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGlamourStyleFromFile(t *testing.T) {
	gotOption := glamour.WithStylesFromJSONBytes(glamourJSONBytes)
	renderer, err := glamour.NewTermRenderer(gotOption)
	require.NoError(t, err)
	assert.NotNil(t, renderer)

	_, err = renderer.Render("a")
	assert.NoError(t, err)
}

func TestGlamourStylesFileIsValid(t *testing.T) {
	got := json.Valid(glamourJSONBytes)
	assert.True(t, got)
}

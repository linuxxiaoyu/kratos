package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTemplateJSON(t *testing.T) {
	rv, err := parseTemplateJSON("./testdata/template.json")
	assert.NoError(t, err)
	fmt.Printf("%v", rv)
}

func TestGenerate(t *testing.T) {
	generateFromTemplateProject("testdata", "/tmp/new_project")
}
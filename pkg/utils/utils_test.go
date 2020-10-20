package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAnotations(t *testing.T) {
	tests := []struct {
		name   string
		expect string
	}{
		{
			name:   "test1",
			expect: "psp-util.test1",
		},
	}

	for _, test := range tests {
		t.Log(test.name, test.expect)
		pspName := GenerateName(test.name)
		ano := GenerateAnotations(pspName)
		val, ok := ano[AnnotaionKeyPSPName]
		assert.True(t, ok)
		assert.Equal(t, test.expect, val)
	}
}

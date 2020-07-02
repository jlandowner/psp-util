package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAnotations(t *testing.T) {
	tests := []struct {
		title   string
		pspName string
		expect  string
	}{
		{
			title:   "':' will be '.'",
			pspName: "psp-util:eks.privileged_01/",
			expect:  "psp-util.eks.privileged_01.",
		},
	}

	for _, test := range tests {
		t.Log(test.title, test.pspName, test.expect)
		ano := GenerateAnotations(test.pspName)
		val, ok := ano[AnnotaionKeyPSPName]
		assert.True(t, ok)
		assert.Equal(t, test.expect, val)
	}
}

package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValuesFromKindFlags(t *testing.T) {
	tests := []struct {
		title       string
		option      AttachDetachOptions
		expectKind  string
		expectName  string
		expectCount int
	}{
		{
			title:       "test1",
			option:      AttachDetachOptions{Group: "testGroup"},
			expectKind:  "Group",
			expectName:  "testGroup",
			expectCount: 1,
		},
		{
			title:       "test2",
			option:      AttachDetachOptions{User: "testUser"},
			expectKind:  "User",
			expectName:  "testUser",
			expectCount: 1,
		},
		{
			title:       "test3",
			option:      AttachDetachOptions{ServiceAccount: "testSA"},
			expectKind:  "ServiceAccount",
			expectName:  "testSA",
			expectCount: 1,
		},
		{
			title:       "test4",
			option:      AttachDetachOptions{Group: "testGroup", User: "testUser"},
			expectKind:  "User",
			expectName:  "testUser",
			expectCount: 2,
		},
		{
			title:       "test5",
			option:      AttachDetachOptions{},
			expectKind:  "",
			expectName:  "",
			expectCount: 0,
		},
	}

	for _, test := range tests {
		t.Log(test.title)
		kind, name, kindFlagCount := getValuesFromKindFlags(&test.option)
		assert.Equal(t, test.expectCount, kindFlagCount)
		assert.Equal(t, test.expectKind, kind)
		assert.Equal(t, test.expectName, name)
	}
}

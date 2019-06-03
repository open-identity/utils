package stringslice

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

func TestDiff(t *testing.T) {

	testData := []struct {
		name            string
		source          []string
		target          []string
		expectedAdded   []string
		expectedRemoved []string
	}{
		{
			name:            "same",
			source:          []string{"1", "2", "3"},
			target:          []string{"1", "2", "3"},
			expectedAdded:   []string{},
			expectedRemoved: []string{},
		},
		{
			name:            "only-source",
			source:          []string{"1", "2", "3"},
			target:          []string{},
			expectedAdded:   []string{},
			expectedRemoved: []string{"1", "2", "3"},
		},
		{
			name:            "only-target",
			source:          []string{},
			target:          []string{"1", "2", "3"},
			expectedAdded:   []string{"1", "2", "3"},
			expectedRemoved: []string{},
		},
		{
			name:            "added-elements",
			source:          []string{"1", "2", "4"},
			target:          []string{"1", "2", "3", "4"},
			expectedAdded:   []string{"3"},
			expectedRemoved: []string{},
		},
		{
			name:            "removed-elements",
			source:          []string{"1", "2", "3", "4"},
			target:          []string{"1", "2", "4"},
			expectedAdded:   []string{},
			expectedRemoved: []string{"3"},
		},
		{
			name:            "mix",
			source:          []string{"1", "10", "2", "3", "4", "7"},
			target:          []string{"1", "2", "4", "5", "6", "7"},
			expectedAdded:   []string{"5", "6"},
			expectedRemoved: []string{"10", "3"},
		},
		{
			name:            "unsorted",
			source:          []string{"3", "4", "1", "2", "7", "10"},
			target:          []string{"5", "6", "7", "1", "2", "4"},
			expectedAdded:   []string{"5", "6"},
			expectedRemoved: []string{"10", "3"},
		},
	}

	for _, td := range testData {
		t.Run(fmt.Sprintf("case=%s", td.name), func(t *testing.T) {
			added, removed := Difference(td.source, td.target)
			assert.Equal(t, td.expectedAdded, added)
			assert.Equal(t, td.expectedRemoved, removed)
		})
	}

}

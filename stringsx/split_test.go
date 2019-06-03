package stringsx_test

import (
	"testing"

	"github.com/open-identity/utils/stringsx"
	"github.com/stretchr/testify/assert"
)

func TestSplitNonEmpty(t *testing.T) {
	// assert.Len(t, strings.Split("", " "), 1)
	assert.Len(t, stringsx.Splitx("", " "), 0)
}

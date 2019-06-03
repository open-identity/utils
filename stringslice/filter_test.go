package stringslice_test

import (
	"testing"
	"unicode"

	"github.com/open-identity/utils/stringslice"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var filter = func(a string) func(b string) bool {
		return func(b string) bool {
			return a == b
		}
	}

	assert.EqualValues(t, []string{"bar"}, stringslice.Filter([]string{"foo", "bar"}, filter("foo")))
	assert.EqualValues(t, []string{"foo"}, stringslice.Filter([]string{"foo", "bar"}, filter("bar")))
	assert.EqualValues(t, []string{"foo", "bar"}, stringslice.Filter([]string{"foo", "bar"}, filter("baz")))
}

func TestTrimEmptyFilter(t *testing.T) {
	assert.EqualValues(t, []string{}, stringslice.TrimEmptyFilter([]string{" ", "  ", "    "}, unicode.IsSpace))
	assert.EqualValues(t, []string{"a"}, stringslice.TrimEmptyFilter([]string{"a", " ", "  ", "    "}, unicode.IsSpace))
}

func TestTrimSpaceEmptyFilter(t *testing.T) {
	assert.EqualValues(t, []string{}, stringslice.TrimSpaceEmptyFilter([]string{" ", "  ", "    "}))
	assert.EqualValues(t, []string{"a"}, stringslice.TrimSpaceEmptyFilter([]string{"a", " ", "  ", "    "}))
}

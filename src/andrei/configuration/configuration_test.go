package configuration

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


var base = "http://google.com/path/to/something"

func TestHostUrlValidatorMatch(t *testing.T) {
	assert.Equal(t, true, HostUrlValidator(base)("http://google.com"))
	assert.Equal(t, true, HostUrlValidator(base)("http://google.com/other/path"))
	assert.Equal(t, true, HostUrlValidator(base)("https://google.com"))
	assert.Equal(t, true, HostUrlValidator(base)("https://google.com/other/path"))
}

func TestHostUrlValidatorNegativeMatch(t *testing.T) {
	assert.Equal(t, false, HostUrlValidator(base)("http://a.b.c"))
	assert.Equal(t, false, HostUrlValidator(base)("https://a.b.c"))
	assert.Equal(t, false, HostUrlValidator(base)("http:/garbage"))
	assert.Equal(t, false, HostUrlValidator(base)("adsfasdfasdf"))
}

func TestAbsolutePathBuilder_whenPathIsPassedInReturnsUrlRelativeToBase(t *testing.T) {
	expected := "http://google.com/path/to/something/else"
	assert.Equal(t, expected, AbsolutePathBuilder(base)("/path/to/something/else"))
}

func TestAbsolutePathBuilder_whenUrlIsPassedInReturnsUrl(t *testing.T) {
	expected := "http://a.b.com"
	assert.Equal(t, expected, AbsolutePathBuilder(base)(expected))
}

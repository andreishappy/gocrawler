package path_utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFunctionalComposition(t *testing.T) {
	assert.Equal(t, true, HostUrlValidator("http://google.com/path/to/something")("http://google.com"))
}


func TestRelativizer_whenPathIsPassedInReturnsUrlRelativeToBase(t *testing.T) {
	relativiser := HostUrlRelativiser("http://google.com/path/to/something")
	assert.Equal(t, "http://google.com/path/to/something/else", relativiser("/path/to/something/else"))
}

func TestRelativizer_whenPathIsPassedInReturnsUrlRelativeToBase(t *testing.T) {
	relativiser := HostUrlRelativiser("http://google.com/path/to/something")
	assert.Equal(t, "http://google2.com/hello", relativiser("http://google2.com/hello"))
}

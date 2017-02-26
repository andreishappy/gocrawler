package path_utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFunctionalComposition(t *testing.T) {
	assert.Equal(t, true, HostUrlValidator("http://google.com/path/to/something")("http://google.com"))
}

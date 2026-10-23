package verdict_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grandper/go-screenplay-udp/internal/verdict"
)

func TestDescribe(t *testing.T) {
	t.Run("quotes a text", func(t *testing.T) {
		assert.Equal(t, `"127.0.0.1:9000"`, verdict.Describe("127.0.0.1:9000"))
	})

	t.Run("quotes bytes", func(t *testing.T) {
		assert.Equal(t, `"Hello"`, verdict.Describe([]byte("Hello")))
	})

	t.Run("writes any other value in its default format", func(t *testing.T) {
		assert.Equal(t, "3", verdict.Describe(3))
		assert.Equal(t, "<nil>", verdict.Describe(nil))
	})
}

package jsonx

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteIndent(t *testing.T) {
	t.Run("indents JSON", func(t *testing.T) {
		v := map[string]interface{}{
			"foo": "bar",
			"baz": []string{"qux"},
		}

		var buf bytes.Buffer

		err := WriteIndent(&buf, v, "", "  ")
		assert.NoError(t, err)

		expected := `{
  "baz": [
    "qux"
  ],
  "foo": "bar"
}`
		assert.Equal(t, expected, buf.String())
	})

	t.Run("returns marshal errors", func(t *testing.T) {
		var buf bytes.Buffer

		err := WriteIndent(&buf, func() {}, "", "  ")
		assert.EqualError(t, err, "json: unsupported type: func()")
	})

	t.Run("returns write errors", func(t *testing.T) {
		err := WriteIndent(badWriter(0), 42, "", "  ")
		assert.EqualError(t, err, "bad writer")
	})
}

type badWriter int

func (badWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("bad writer")
}

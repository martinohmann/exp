package json

import (
	"encoding/json"
	"io"
)

// WriteIndent is like encoding/json.MarshalIndent but instead of returning the
// indented JSON bytes for v, it writes them to w.
//
// See stdlib encoding/json.MarshalIndent for more information.
func WriteIndent(w io.Writer, v interface{}, prefix, indent string) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	return err
}

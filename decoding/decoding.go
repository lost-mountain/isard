package decoding

import (
	"encoding/json"
	"io"
)

// Decoder is an interface to decode incoming
// messages.
type Decoder func(o interface{}, r io.Reader) error

// JSONDecoder implements the Decoder interface
// to decode messages in JSON representation.
func JSONDecoder(o interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(o)
}

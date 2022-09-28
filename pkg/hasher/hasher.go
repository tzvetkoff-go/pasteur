package hasher

import (
	"github.com/speps/go-hashids"
	"github.com/tzvetkoff-go/errors"
)

// Hasher ...
type Hasher struct {
	Hashids *hashids.HashID
}

// New ...
func New(config *Config) (*Hasher, error) {
	hd := hashids.NewData()
	hd.Alphabet = config.Alphabet
	hd.Salt = config.Salt
	hd.MinLength = config.MinLength

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, errors.Propagate(err, "cannot create hasher")
	}

	result := &Hasher{
		Hashids: h,
	}
	return result, nil
}

// Encode ...
func (h *Hasher) Encode(id int) (string, error) {
	return h.Hashids.Encode([]int{id})
}

// Decode ...
func (h *Hasher) Decode(hash string) (int, error) {
	ints, err := h.Hashids.DecodeWithError(hash)
	if err != nil {
		return 0, errors.Propagate(err, "cannot decode hash")
	}
	if len(ints) != 1 {
		return 0, errors.Propagate(err, "expected exactly 1 number in hash, got %d", len(ints))
	}

	return ints[0], nil
}

package suggester

import "errors"

var (
	// ErrNoMatch は BaseText 内に置換対象が見つからなかったことを示す。
	ErrNoMatch = errors.New("no match location found")
	// ErrEmptyAfter は LLMAfter が空だったことを示す。
	ErrEmptyAfter = errors.New("empty LLMAfter")
	// ErrNoChange は BaseText と LLMAfter の差分が無いことを示す。
	ErrNoChange = errors.New("no change between BaseText and LLMAfter")
)

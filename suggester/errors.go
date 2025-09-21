package suggester

import "errors"

var (
	ErrNoMatch    = errors.New("no match location found")
	ErrEmptyAfter = errors.New("empty LLMAfter")
	ErrNoChange   = errors.New("no change between BaseText and LLMAfter")
)
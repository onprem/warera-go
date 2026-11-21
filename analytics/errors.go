package analytics

import "errors"

// Errors the analytics packages return for codes and levels the game
// configuration does not define. Test for them with [errors.Is].
var (
	ErrUnknownItem    = errors.New("analytics: unknown item")
	ErrUnknownSkill   = errors.New("analytics: unknown skill")
	ErrUnknownUpgrade = errors.New("analytics: unknown upgrade")
)

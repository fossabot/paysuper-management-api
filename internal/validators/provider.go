package validators

import (
	"github.com/google/wire"
)

// Provider
func Provider() (*ValidatorSet, func(), error) {
	g := New()
	return g, func() {}, nil
}

var (
	WireSet = wire.NewSet(
		Provider,
	)
	WireTestSet = wire.NewSet(
		Provider,
	)
)

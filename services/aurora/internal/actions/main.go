package actions

import "github.com/diamcircle/go/services/aurora/internal/corestate"

type CoreStateGetter interface {
	GetCoreState() corestate.State
}

package emitter

import (
	"github.com/K-Phoen/flatrise/flatrise/model"
)

type Emitter interface {
	Emit(offer model.Offer) error
}

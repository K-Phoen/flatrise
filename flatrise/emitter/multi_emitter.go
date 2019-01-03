package emitter

import (
	"github.com/K-Phoen/flatrise/flatrise/model"
)

type multiEmitter struct {
	emitters []Emitter
}

func NewMultiEmitter(emitters ...Emitter) Emitter {
	return multiEmitter{
		emitters: emitters,
	}
}

func (emitter multiEmitter) Emit(offer model.Offer) error {
	for _, e := range emitter.emitters {
		err := e.Emit(offer)
		if err != nil {
			return err
		}
	}

	return nil
}

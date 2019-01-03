package emitter

import (
	"encoding/json"
	"github.com/K-Phoen/flatrise/flatrise/model"
	"github.com/pkg/errors"
	"io"
)

type writerEmitter struct {
	writer io.Writer
}

func NewWriterEmitter(writer io.Writer) writerEmitter {
	return writerEmitter{
		writer: writer,
	}
}

func (emitter writerEmitter) Emit(offer model.Offer) error {
	jsonOffer, err := json.Marshal(offer)
	if err != nil {
		return errors.Wrap(err, "Could not marshal offer")
	}

	jsonOffer = append(jsonOffer, '\n')
	_, err = emitter.writer.Write(jsonOffer)

	return err
}

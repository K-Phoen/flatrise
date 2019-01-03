package worker

import (
	"github.com/K-Phoen/flatrise/flatrise/emitter"
)

type CrawlerFunc func(emitter emitter.Emitter)

type Worker interface {
	Run() error
}

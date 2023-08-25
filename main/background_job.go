package main

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	fn func() error
}

func NewWorker(fn func() error) *Worker {
	return &Worker{fn: fn}
}

func (w *Worker) RunBackground(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	w.tick()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick()
		}
	}
}

func (w *Worker) tick() {
	if err := w.fn(); err != nil {
		log.Print(err)
	}
}

package phudong

import (
	"context"
	"sync"
	"time"
)

type Worker struct {
	waitFunc func()

	options
}

func NewWorker(opts ...optionFunc) *Worker {
	return &Worker{
		options: newOptions(opts...),
	}
}

func (w *Worker) Wait() {
	if w.waitFunc != nil {
		w.waitFunc()
	}
}

func (w *Worker) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	w.waitFunc = wg.Wait

	ticker := time.NewTicker(w.duration)

	processError := func(ctx context.Context, err error) {
		if w.withErrorProcessor != nil {
			w.withErrorProcessor(ctx, err)
		}
	}

	doThisWrapper := func(ctx context.Context) {
		if w.doThis == nil && w.doThisOrThrowError == nil {
			w.logger.Errorf("%s: no function set to execute", w.name)
			processError(ctx, ErrNoFunctionSet)
		}

		if w.doThis != nil {
			w.doThis(ctx)
		}

		if w.doThisOrThrowError != nil {
			if err := w.doThisOrThrowError(ctx); err != nil {
				w.logger.Errorf("%s: error executing function: %v", w.name, err)
				processError(ctx, err)
			}
		}
	}

	go func() {
		w.logger.Debugf(w.name + " started")
		defer w.logger.Debugf(w.name + " stopped")
		defer wg.Done()
		defer ticker.Stop()

		if w.instantRun {
			doThisWrapper(ctx)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				doThisWrapper(ctx)
			}
		}
	}()
}

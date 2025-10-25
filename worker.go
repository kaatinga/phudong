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
		if len(w.doThis) == 0 && len(w.doThisOrThrowError) == 0 {
			w.logger.Errorf("%s: no function set to execute\n", w.name)
			processError(ctx, ErrNoFunctionSet)
		}

		for _, do := range w.doThis {
			do(ctx)
		}

		for _, do := range w.doThisOrThrowError {
			if err := do(ctx); err != nil {
				w.logger.Errorf("%s: error executing function: %v\n", w.name, err)
				processError(ctx, err)
			}
		}
	}

	go func() {
		w.logger.Printf(w.name + " started\n")
		defer w.logger.Printf(w.name + " stopped\n")
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

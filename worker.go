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

func (w *Worker) currentTime() time.Time {
	if w.now != nil {
		return w.now()
	}
	return time.Now()
}

func (w *Worker) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	w.waitFunc = wg.Wait

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
		w.logger.Printf(w.name + " started")
		defer w.logger.Printf(w.name + " stopped")
		defer wg.Done()

		if w.instantRun {
			doThisWrapper(ctx)
		}

		if w.dailyAt {
			w.runDaily(ctx, doThisWrapper)
			return
		}

		ticker := time.NewTicker(w.duration)
		defer ticker.Stop()

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

func (w *Worker) runDaily(ctx context.Context, do func(context.Context)) {
	for {
		now := w.currentTime()
		delay := max(nextDailyAt(now, w.dailyHour, w.dailyMinute).Sub(now), 0)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		case <-timer.C:
			do(ctx)
		}
	}
}

// nextDailyAt returns the next UTC occurrence of hour:minute.
// If that time has already been reached today, the result is tomorrow.
func nextDailyAt(now time.Time, hour, minute byte) time.Time {
	now = now.UTC()
	next := time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(minute), 0, 0, time.UTC)
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

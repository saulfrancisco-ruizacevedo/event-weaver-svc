package common

import (
	"context"
	"sync"
)

func FetchWithChannel[T any](ctx context.Context, repoCall func(context.Context) (T, error), resultCh chan<- T, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := repoCall(ctx)

	if err != nil {
		errCh <- err
		return
	}

	resultCh <- res
}

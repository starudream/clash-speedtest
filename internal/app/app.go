package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

var (
	signals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT}

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	eg     *errgroup.Group

	ss []S
)

type S func(ctx context.Context) error

func Add(s S) {
	mu.Lock()
	defer mu.Unlock()
	ss = append(ss, s)
}

func Go() error {
	return internalGo(false)
}

func OnceGo() error {
	return internalGo(true)
}

func internalGo(once bool) error {
	mu.Lock()
	defer mu.Unlock()

	if len(ss) == 0 {
		return nil
	}

	ctx, cancel = context.WithCancel(context.Background())
	eg, ctx = errgroup.WithContext(ctx)

	errCh := make(chan error, 1)
	for i := 0; i < len(ss); i++ {
		s := ss[i]
		if once {
			eg.Go(func() error { return s(ctx) })
		} else {
			go func() { errCh <- s(ctx) }()
		}
	}

	var err error

	if !once {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, signals...)

		eg.Go(func() error {
			for {
				select {
				case err = <-errCh:
					if err != nil {
						Stop()
					}
				case <-ch:
					Stop()
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	ege := eg.Wait()

	if err != nil {
		return err
	}

	if ege != nil {
		if !errors.Is(ege, context.Canceled) {
			err = ege
		}
		return err
	}

	return nil
}

func Stop() {
	if cancel != nil {
		fmt.Println()
		cancel()
	}
}

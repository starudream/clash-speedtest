package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
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
	fs []F
)

type S func(ctx context.Context) error

func Add(s S) {
	mu.Lock()
	defer mu.Unlock()
	ss = append(ss, s)
}

type F func()

func Recover(_fs ...F) {
	mu.Lock()
	defer mu.Unlock()
	fs = append(fs, _fs...)
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

	wg := sync.WaitGroup{}
	wg.Add(len(ss))

	doneCh := make(chan struct{}, 1)

	if once {
		go func() {
			wg.Wait()
			close(doneCh)
		}()
	}

	errCh := make(chan error, 1)

	for i := 0; i < len(ss); i++ {
		s := ss[i]
		go func() {
			defer func() {
				wg.Done()
				if re := recover(); re != nil {
					fmt.Println(string(debug.Stack()))
					Stop()
				}
			}()
			e := s(ctx)
			if e != nil {
				errCh <- e
			}
		}()
	}

	var err error

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, signals...)

	eg.Go(func() error {
		for {
			select {
			case err = <-errCh:
				Stop()
			case <-doneCh:
				Stop()
			case <-signalCh:
				Stop()
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	ege := eg.Wait()

	for i := len(fs) - 1; i >= 0; i-- {
		fs[i]()
	}

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

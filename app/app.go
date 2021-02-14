package app

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"moocss.com/tiga/pkg/log"
	"os"
	"os/signal"
)

// App .
type App struct {
	opts   *options
	log    *log.Helper
	ctx    context.Context
	cancel func()
}

// New .
func NewApp(opts ...Option) *App {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
		log:    log.NewHelper("app", options.logger),
	}
}

// Run .
func (a *App) Run() error {
	g, ctx := errgroup.WithContext(a.ctx)

	srv := a.opts.server

	//在服务启动前执行钩子
	//srv.AddBeforeServerStartFunc(
	//	//...
	//)
	//
	//在服务关闭前执行钩子
	//srv.AddAfterServerStopFunc(
	//	//...
	//)

	g.Go(func() error {
		<-ctx.Done() // wait for stop signal
		return srv.Stop()
	})
	g.Go(func() error {
		// return srv.Run()
		return srv.Start()
	})

	// signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, a.opts.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				a.Stop()
			}
		}
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

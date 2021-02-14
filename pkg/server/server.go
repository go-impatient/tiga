package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/server/bootstrap"
)

// Server is a HTTP server.
type Server struct {
	*http.Server
	beforeFuncs []bootstrap.BeforeServerStartFunc
	afterFuncs  []bootstrap.AfterServerStopFunc
	opts        *options
	log         *log.Helper
}

// NewServer new a HTTP server.
func NewServer(opts ...Option) *Server {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}
	srv := &Server{
		opts: options,
		log:  log.NewHelper("http", options.logger),
	}

	srv.Server = &http.Server{Handler: srv.opts.handler}

	if srv.opts.middleware != nil {
		srv.Handler = srv.opts.middleware(srv.opts.handler)
	}

	return srv
}

// Serve serve http request
func (s *Server) Run() (err error) {
	for _, fn := range s.beforeFuncs {
		err = fn()
		if err != nil {
			return err
		}
	}

	s.Start()

	for _, fn := range s.afterFuncs {
		fn()
	}

	return err
}

func (s *Server) Stop() error {
	s.log.Info("[HTTP] server stopping")
	return s.Shutdown(context.Background())
}

func (s *Server) Start() error {
	lis, err := net.Listen(s.opts.network, s.opts.address)
	if err != nil {
		return err
	}
	s.log.Infof("[HTTP] server listening on: %s", s.opts.address)
	return s.Serve(lis)
}

// PingServer 服务心跳检查
func (s *Server) PingServer() error {
	maxPingCount := conf.GetInt("app.max_ping_count")
	if maxPingCount == 0 {
		maxPingCount = 2
	}
	for i := 0; i < maxPingCount; i++ {
		// Ping the app by sending a GET request to `/health`.
		url := fmt.Sprintf("%s/sd/health", s.opts.address)
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		s.log.Warnw("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}

func (s *Server) redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func (s *Server) AddBeforeServerStartFunc(fns ...bootstrap.BeforeServerStartFunc) {
	for _, fn := range fns {
		s.beforeFuncs = append(s.beforeFuncs, fn)
	}
}

func (s *Server) AddAfterServerStopFunc(fns ...bootstrap.AfterServerStopFunc) {
	for _, fn := range fns {
		s.afterFuncs = append(s.afterFuncs, fn)
	}
}
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	appName    string
	env        string
	port       int
	err        error
	httpServer *http.Server
	mux        *http.ServeMux
}

type Option func(*Server)

func WithAppName(appName string) Option {
	return func(s *Server) {
		s.appName = appName
	}
}

func WithEnv(env string) Option {
	return func(s *Server) {
		s.env = env
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func New(opts ...Option) *Server {
	s := &Server{
		port: 8080,
	}

	for _, op := range opts {
		op(s)
	}

	mux := http.NewServeMux()
	s.mux = mux
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	return s
}

func (s *Server) Run() {
	log.Printf("[%s:%s] server started on :%d", s.appName, s.env, s.port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %+v", err)
		}
	}()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server exit")
}

func (s *Server) GET(path string, handler HandlerFunc) {
	s.handle(http.MethodGet, path, handler)
}

func (s *Server) POST(path string, handler HandlerFunc) {
	s.handle(http.MethodPost, path, handler)
}

func (s *Server) handle(method string, path string, handler HandlerFunc) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
		}

		ctx := &Context{
			Writer:  w,
			Request: r,
		}
		handler(ctx)
	})
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

type HandlerFunc func(*Context)

func (c *Context) JSON(status int, data any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		http.Error(c.Writer, "internal error", http.StatusInternalServerError)
	}
}

func (c *Context) OK(data any) {
	c.JSON(http.StatusOK, data)
}

func (c *Context) InternalServerError(data any) {
	c.JSON(http.StatusInternalServerError, data)
}

func (c *Context) BindJSON(v any) error {
	if c.Request.Body == nil {
		return fmt.Errorf("empty request body")
	}

	defer c.Request.Body.Close()

	if err := json.NewDecoder(c.Request.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

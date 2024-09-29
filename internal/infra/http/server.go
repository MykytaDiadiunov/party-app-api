package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Server(ctx context.Context, router http.Handler) error {
	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", 8080),
	}

	errServeCh := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errServeCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error during server shutdown: %w", err)
		}
	case err := <-errServeCh:
		return fmt.Errorf("error during server execution: %w", err)
	}
	return nil
}

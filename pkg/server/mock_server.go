package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	m "github.com/bmaynard/apimock/pkg/mocks"

	"github.com/bmaynard/apimock/pkg/responses"
	l "github.com/bmaynard/apimock/pkg/utils/logger"
)

var mocks map[string]m.MockResponse

func Serve(wait time.Duration, addr string, pemPath string, keyPath string) {
	mocks = make(map[string]m.MockResponse)

	r := mux.NewRouter().StrictSlash(true)
	responses.BuildRoutes(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if len(pemPath) > 1 && len(keyPath) > 1 {
			if err := srv.ListenAndServeTLS(pemPath, keyPath); err != nil {
				l.Log.Fatal(err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil {
				l.Log.Fatal(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	l.Log.Info("shutting down")
	os.Exit(0)
}

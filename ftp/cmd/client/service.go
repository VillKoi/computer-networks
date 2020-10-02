package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/VillKoi/computer-networks/ftp/client"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"golang.org/x/sync/errgroup"
)

func StartHTTP(ctx context.Context, c *client.FTPClient) error {
	router := chi.NewRouter()

	router.Post("/auth", c.Auth)

	router.Get("/ls", c.Ls)

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1234"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	})

	srv := http.Server{
		Addr:    httpport,
		Handler: cor.Handler(router),
	}

	group := errgroup.Group{}
	group.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		return srv.Shutdown(ctx)
	})

	return group.Wait()
}

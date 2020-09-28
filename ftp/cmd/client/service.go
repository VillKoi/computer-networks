package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/VillKoi/computer-networks/ftp/client"
	"golang.org/x/sync/errgroup"
)

func StartHTTP(ctx context.Context, c *client.FTPClient) {

	mux := http.NewServeMux()

	mux.HandleFunc("/auth", c.Auth)

	var handler http.Handler = mux

	// for _, m := range middlewares {
	// 	handler = m(handler)
	// }

	srv := http.Server{
		Addr:    httpport,
		Handler: handler,
	}

	group := errgroup.Group{}
	group.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
}

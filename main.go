package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/axyut/dairygo/internal/config"
	"github.com/axyut/dairygo/internal/db"
	"github.com/axyut/dairygo/internal/handler"
	"github.com/axyut/dairygo/internal/service"
)

var ctx = context.Background()

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	conf := config.Envs
	mongo, err := db.NewMongo(ctx, conf, logger)
	if err != nil {
		panic(err)
	}
	srv := service.NewService(ctx, mongo, logger)
	server := handler.RootHandler(ctx, conf, srv, logger)

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	go func() {
		server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Log(ctx, slog.LevelInfo,
		"Connected to Server",
		slog.String("port", conf.PORT),
		slog.String("env", conf.ENV),
		slog.String("server", "http://localhost:"+conf.PORT),
	)

	<-killSig

	// Server shutdown
	if err := mongo.Close(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}
	// cancel()
	logger.Info("Server shutdown complete")
}

// buy at different rate from same vendor for different goods
// when transaction deleted, update audience to pay and to receive
// when transaction updated, update audience to pay and to receive
// when transaction inserted(bought,sold), update audience to pay and to receive

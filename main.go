package main

import (
	"errors"
	"fmt"
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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	conf := config.Envs
	mongo, err := db.NewMongo(conf, logger)
	if err != nil {
		panic(err)
	}
	srv := service.NewService(mongo, logger)
	handler.RootHandler(srv, logger)

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%s", conf.PORT), nil)

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Log(db.Ctx, slog.LevelInfo,
		"Connected to Server",
		slog.String("port", conf.PORT),
		slog.String("env", conf.ENV),
		slog.String("server", "http://localhost:"+conf.PORT),
	)

	<-killSig

	// Server shutdown
	if err := mongo.Close(db.Ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}
	logger.Info("Server shutdown complete")
}

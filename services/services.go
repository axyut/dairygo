package services

import (
	"github.com/axyut/dairygo/db"

	"golang.org/x/exp/slog"
)

type Service struct {
	Log *slog.Logger
	db  *db.Mongo
}

func NewService(log *slog.Logger, db *db.Mongo) *Service {
	return &Service{log, db}
}

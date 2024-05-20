package services

import (
	"github.com/axyut/dairygo/db"

	"golang.org/x/exp/slog"
)

func GlobalService(log *slog.Logger, cs *db.MyDb) string {
	return "Global Service Usage"
}

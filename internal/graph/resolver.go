package graph

import (
	"log/slog"

	"github.com/augustdev/autoclip/internal/storage/pg"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Db     *pg.DB
	Logger *slog.Logger
}

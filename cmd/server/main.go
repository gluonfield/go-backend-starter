package main

import (
	"github.com/augustdev/autoclip/internal/bootstrap"
	"github.com/augustdev/autoclip/internal/storage/pg"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			bootstrap.NewLogger,
			bootstrap.NewConfig,
			pg.NewDatabase,
			bootstrap.NewResolver,
			bootstrap.NewTokenValidator,
			bootstrap.NewGraphQLRouter,
		),
		fx.Invoke(
			bootstrap.StartServer,
		),
	).Run()
}

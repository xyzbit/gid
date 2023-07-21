//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/xyzbit/gid/core/conf"
	"github.com/xyzbit/gid/core/loader"
	"github.com/xyzbit/gid/core/server"
)

// initApp init application.
func initGrpcServer(context.Context, *conf.Server, *conf.DBConfig, *conf.ConsulConfig) (*server.GrpcServer, error) {
	panic(wire.Build(server.ProviderSet, loader.ProviderSet))
}

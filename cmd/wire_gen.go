// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/xyzbit/gid/core/conf"
	"github.com/xyzbit/gid/core/loader"
	"github.com/xyzbit/gid/core/server"
)

// Injectors from wire.go:

// initApp init application.
func initGrpcServer(contextContext context.Context, confServer *conf.Server, dbConfig *conf.DBConfig, consulConfig *conf.ConsulConfig) (*server.GrpcServer, error) {
	statusLoader := loader.NewStatusLoader(dbConfig)
	configLoader := loader.NewConfigLoader(consulConfig)
	generatorSvc := server.NewGeneratorSvc(statusLoader, configLoader)
	grpcServer := server.NewGrpcServer(generatorSvc, confServer)
	return grpcServer, nil
}

package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	v1 "github.com/xyzbit/gid/api/v1"
	"github.com/xyzbit/gid/core"
)

type GeneratorSvc struct {
	cl       core.ConfigLoader
	sl       core.StatusLoader
	modulesG map[string]*core.IDGenerator
	v1.UnimplementedGeneratorServer
	v1.UnimplementedMannagerServer
}

func NewGeneratorSvc(sl core.StatusLoader, cl core.ConfigLoader) *GeneratorSvc {
	configs, err := cl.AllConfigs()
	if err != nil {
		panic(err)
	}
	modulesG := make(map[string]*core.IDGenerator, len(configs))
	for _, cfg := range configs {
		modulesG[cfg.ModuleName] = core.NewGenerator(cfg.ModuleName, sl, core.WithConfig(cfg))
	}
	return &GeneratorSvc{
		cl:       cl,
		sl:       sl,
		modulesG: modulesG,
	}
}

func (gsvc *GeneratorSvc) NextID(ctx context.Context, req *v1.NextIDReq) (*v1.NextIDReply, error) {
	if req.Module == "" {
		return nil, errors.New("module name is empty")
	}
	g, ok := gsvc.modulesG[req.Module]
	if !ok {
		return nil, errors.New("module not found")
	}

	id := g.NextID()
	if id == -1 {
		return nil, errors.Errorf("id use over! please check that the id scheduler is healthy")
	}

	return &v1.NextIDReply{
		Id: id,
	}, nil
}

func (gsvc *GeneratorSvc) Moudles(ctx context.Context, req *v1.MoudlesReq) (*v1.MoudlesReply, error) {
	configs, err := gsvc.cl.AllConfigs()
	if err != nil {
		return nil, err
	}
	status, err := gsvc.sl.AllStatus()
	if err != nil {
		return nil, err
	}
	moudles := make([]*v1.Moudle, len(configs))
	for i, cfg := range configs {
		var s *core.Status
		for _, item := range status {
			if item.ModuleName == cfg.ModuleName {
				s = item
			}
		}
		moudles[i] = &v1.Moudle{
			Name:          cfg.ModuleName,
			Desc:          cfg.Desc,
			IdKind:        v1.IDKind(cfg.IDKind),
			FirstId:       cfg.IncrFirstID,
			LoadNum:       cfg.LoadNum,
			LoadThreshold: cfg.LoadThreshold,
		}
		if s != nil {
			moudles[i].LastId = s.LastID
		}
	}
	return &v1.MoudlesReply{
		Moudles: moudles,
	}, nil
}

func (gsvc *GeneratorSvc) RegisterModules(ctx context.Context, req *v1.RegisterModulesReq) (*v1.RegisterModulesReply, error) {
	if req.Name == "" {
		return nil, errors.New("module name is empty")
	}
	if _, ok := gsvc.modulesG[req.Name]; ok {
		return nil, fmt.Errorf("module %s already exists", req.Name)
	}

	switch req.IdKind {
	case v1.IDKind_IncrementID:
		option := req.GetIncrementOption()
		cfg := &core.Config{
			ModuleName:    req.Name,
			Desc:          req.Desc,
			IDKind:        int8(req.IdKind),
			LoadNum:       option.LoadNum,
			LoadThreshold: option.LoadThreshold,
			IncrFirstID:   option.FirstId,
		}
		if err := gsvc.cl.SaveConfig(req.Name, cfg); err != nil {
			return nil, err
		}
		gsvc.modulesG[req.Name] = core.NewGenerator(cfg.ModuleName, gsvc.sl, core.WithConfig(cfg))
	case v1.IDKind_SnowflakeID:
		return nil, errors.New("snowflake id not implemented")
	}
	return &v1.RegisterModulesReply{}, nil
}

// TODO: unregister module
// func (gsvc *GeneratorSvc) UnregisterModules(ctx context.Context, req *v1.UnregisterModulesReq) (*v1.UnregisterModulesReply, error) {
// }

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/29

package gomodule

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"time"
)

type Module interface {
	VersionList(ctx context.Context) ([]string, error)
	VersionInfo(ctx context.Context, version string) (*Info, error)
	VersionMod(ctx context.Context, version string) ([]byte, error)
	VersionFiles(ctx context.Context, version string) ([]fs.DirEntry, error)
	Latest(ctx context.Context) (*Info, error)
}

type Info struct {
	// Time commit time
	Time time.Time

	// Version
	Version string
}

func (in *Info) String() string {
	bf, _ := json.Marshal(in)
	return string(bf)
}

var errEmptyModules = errors.New("empty modules")

var _ Module = (Modules)(nil)

type Modules []Module

func (ms Modules) VersionList(ctx context.Context) ([]string, error) {
	if len(ms) == 0 {
		return nil, errEmptyModules
	}
	var err error
	for _, m := range ms {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var ret []string
		ret, err = m.VersionList(ctx)
		if err == nil {
			return ret, err
		}
	}
	return nil, err
}

func (ms Modules) VersionInfo(ctx context.Context, version string) (*Info, error) {
	if len(ms) == 0 {
		return nil, errEmptyModules
	}
	var err error
	for _, m := range ms {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var ret *Info
		ret, err = m.VersionInfo(ctx, version)
		if err == nil {
			return ret, err
		}
	}
	return nil, err
}

func (ms Modules) VersionMod(ctx context.Context, version string) ([]byte, error) {
	if len(ms) == 0 {
		return nil, errEmptyModules
	}
	var err error
	for _, m := range ms {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var ret []byte
		ret, err = m.VersionMod(ctx, version)
		if err == nil {
			return ret, err
		}
	}
	return nil, err
}

func (ms Modules) VersionFiles(ctx context.Context, version string) ([]fs.DirEntry, error) {
	if len(ms) == 0 {
		return nil, errEmptyModules
	}
	var err error
	for _, m := range ms {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var ret []fs.DirEntry
		ret, err = m.VersionFiles(ctx, version)
		if err == nil {
			return ret, err
		}
	}
	return nil, err
}

func (ms Modules) Latest(ctx context.Context) (*Info, error) {
	if len(ms) == 0 {
		return nil, errEmptyModules
	}
	var err error
	for _, m := range ms {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		var ret *Info
		ret, err = m.Latest(ctx)
		if err == nil {
			return ret, err
		}
	}
	return nil, err
}

func (ms *Modules) Append(items ...Module) {
	*ms = append(*ms, items...)
}

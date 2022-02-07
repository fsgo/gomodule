// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/29

package gomodule

import (
	"context"
	"encoding/json"
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
	// Version
	Version string

	// Time commit time
	Time time.Time
}

func (in *Info) String() string {
	bf, _ := json.Marshal(in)
	return string(bf)
}

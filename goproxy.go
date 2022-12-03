// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/30

package gomodule

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsgo/cmdutil/gosdk"
)

var defaultUA = "fsgo/gomodule"

var _ Module = (*GoProxy)(nil)

// GoProxy 通过 goproxy 获取 模块的信息，如版本列表，指定版本号的元信息、代码
//
// https://go.dev/ref/mod#module-proxy
type GoProxy struct {
	Client HTTPClient
	Proxy  string
	Module string
}

func (m *GoProxy) client() HTTPClient {
	return client(m.Client)
}

func (m *GoProxy) getProxy() string {
	if len(m.Proxy) > 0 {
		return m.Proxy
	}
	return goProxyFromEnv()
}

const defaultGoProxy = "https://goproxy.cn/"

func goProxyFromEnv() string {
	ev := os.Getenv("GOPROXY")
	if len(ev) == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		goBin := gosdk.LatestOrDefault()
		cmd := exec.CommandContext(ctx, goBin, "env", "GOPROXY")
		bs, err := cmd.Output()
		if err == nil {
			ev = string(bs)
		}
	}
	ev = strings.TrimSpace(ev)
	if len(ev) == 0 {
		return defaultGoProxy
	}
	ss := strings.Split(ev, ",")
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "http") {
			return s
		}
	}
	return defaultGoProxy
}

func (m *GoProxy) query(ctx context.Context, api string) ([]byte, error) {
	proxy := m.getProxy()
	if len(proxy) == 0 {
		return nil, errors.New("empty GoProxy")
	}
	var b strings.Builder
	b.WriteString(proxy)
	if !strings.HasSuffix(m.Proxy, "/") {
		b.WriteString("/")
	}
	b.WriteString(m.Module)
	b.WriteString(api)

	return sentRequest(ctx, m.client(), http.MethodGet, b.String())
}

// VersionList 版本列表
// call $base/$module/@v/list
func (m *GoProxy) VersionList(ctx context.Context) ([]string, error) {
	bf, err := m.query(ctx, "/@v/list")
	if err != nil {
		return nil, err
	}
	bf = bytes.TrimSpace(bf)
	return strings.Split(string(bf), "\n"), nil
}

func (m *GoProxy) queryInfo(ctx context.Context, api string) (*Info, error) {
	bf, err := m.query(ctx, api)
	if err != nil {
		return nil, err
	}
	var info *Info
	err = json.Unmarshal(bf, &info)
	return info, err
}

// VersionInfo
//
// call $base/$module/@v/$version.info
func (m *GoProxy) VersionInfo(ctx context.Context, version string) (*Info, error) {
	return m.queryInfo(ctx, "/@v/"+version+".info")
}

// VersionMod go.mod content
//
// call $base/$module/@v/$version.mod
func (m *GoProxy) VersionMod(ctx context.Context, version string) ([]byte, error) {
	return m.query(ctx, "/@v/"+version+".mod")
}

// VersionZip 返回原始的 zip 数据
// $base/$module/@v/$version.zip
func (m *GoProxy) VersionZip(ctx context.Context, version string) (*zip.Reader, error) {
	bf, err := m.query(ctx, "/@v/"+version+".zip")
	if err != nil {
		return nil, err
	}
	bb := bytes.NewReader(bf)
	return zip.NewReader(bb, int64(len(bf)))
}

func (m *GoProxy) VersionFiles(ctx context.Context, version string) ([]fs.DirEntry, error) {
	zr, err := m.VersionZip(ctx, version)
	if err != nil {
		return nil, err
	}
	prefix := fmt.Sprintf("%s@%s/", m.Module, version)

	var result []fs.DirEntry
	for _, zf := range zr.File {
		name := zf.Name[len(prefix):]
		if len(name) == 0 {
			continue
		}
		isDir := strings.HasSuffix(name, "/")
		if isDir {
			name = name[:len(name)-1]
		}
		zd := &zipDirEntry{
			name: name,
			zf:   zf,
		}
		result = append(result, zd)
	}

	return result, nil
}

// Latest version
//
// call $base/$module/@latest
func (m *GoProxy) Latest(ctx context.Context) (*Info, error) {
	return m.queryInfo(ctx, "/@latest")
}

var _ fs.DirEntry = (*zipDirEntry)(nil)

type zipDirEntry struct {
	zf   *zip.File
	name string
}

func (z *zipDirEntry) Name() string {
	return z.name
}

func (z *zipDirEntry) IsDir() bool {
	return z.zf.FileInfo().IsDir()
}

func (z *zipDirEntry) Type() fs.FileMode {
	return z.zf.Mode()
}

func (z *zipDirEntry) Info() (fs.FileInfo, error) {
	return z.zf.FileInfo(), nil
}

// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/3

package gomodule

import (
	"bytes"
	"context"
	"os/exec"
	"sync"
	"time"

	"github.com/fsgo/cmdutil/gosdk"
	"golang.org/x/mod/module"
)

var (
	goNoProxy     string
	goNoProxyErr  error
	goNoProxyOnce sync.Once
)

// IsNoGoProxy 判断地址是否不使用 GoProxy
func IsNoGoProxy(path string) (bool, error) {
	goNoProxyOnce.Do(func() {
		goNoProxy, goNoProxyErr = loadNoProxy()
	})
	if goNoProxyErr != nil {
		return false, goNoProxyErr
	}
	return module.MatchPrefixPatterns(goNoProxy, path), nil
}

func loadNoProxy() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	goBin := gosdk.LatestOrDefault()
	cmd := exec.CommandContext(ctx, goBin, "env", "GONOPROXY")
	bs, err := cmd.Output()
	if err != nil {
		return "", err
	}
	bs = bytes.TrimSpace(bs)
	return string(bs), nil
}

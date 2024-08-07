// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/3

package gomodule

import (
	"sync"
	"testing"

	"github.com/fsgo/fst"
)

func TestIsNoGoProxy(t *testing.T) {
	goNoProxyOnce = sync.Once{}
	path := "github.com/fsgo/gomodule"
	t.Run("case 1", func(t *testing.T) {
		t.Setenv("GONOPROXY", "other")
		got, e1 := IsNoGoProxy(path)
		fst.NoError(t, e1)
		fst.False(t, got)

		got, e1 = IsNoGoProxy("other/hello")
		fst.NoError(t, e1)
		fst.True(t, got)
	})

	t.Run("case 1", func(t *testing.T) {
		goNoProxyOnce = sync.Once{}
		t.Setenv("GONOPROXY", "github.com*,example.com*")
		got, e1 := IsNoGoProxy(path)
		fst.NoError(t, e1)
		fst.True(t, got)

		got, e1 = IsNoGoProxy("other/hello")
		fst.NoError(t, e1)
		fst.False(t, got)

		got, e1 = IsNoGoProxy("example.com/abc")
		fst.NoError(t, e1)
		fst.True(t, got)
	})
}

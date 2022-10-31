// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/10/31

package gomodule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_goProxyFromEnv(t *testing.T) {
	t.Run("invalid env", func(t *testing.T) {
		t.Setenv("GOPROXY", "abcd")
		require.Equal(t, defaultGoProxy, goProxyFromEnv())
	})
	t.Run("with env", func(t *testing.T) {
		t.Setenv("GOPROXY", "http://abcd")
		require.Equal(t, "http://abcd", goProxyFromEnv())
	})
}

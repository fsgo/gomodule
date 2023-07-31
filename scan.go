// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/7/31

package gomodule

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// ScanGoModFile 查找 go.mod 文件
func ScanGoModFile(fromDir string, fn func(dir string, mod modfile.File) error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	scan := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if err = os.Chdir(wd); err != nil {
			return err
		}
		name := filepath.Base(path)
		if name != "go.mod" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := modfile.Parse("go.mod", content, nil)
		if err != nil {
			return err
		}
		return fn(filepath.Dir(path), *f)
	}
	return filepath.Walk(fromDir, scan)
}

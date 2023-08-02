// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/7/31

package gomodule

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

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

func ScanGoModFileParallel(fromDir string, conc int, fn func(dir string, mod modfile.File) error) error {
	if conc <= 1 {
		return ScanGoModFile(fromDir, fn)
	}
	limiter := make(chan struct{}, conc)
	var wg sync.WaitGroup
	var fnErr error
	var mux sync.Mutex
	err := ScanGoModFile(fromDir, func(dir string, mod modfile.File) error {
		limiter <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-limiter
				wg.Done()
			}()
			if err1 := fn(dir, mod); err1 != nil {
				mux.Lock()
				fnErr = err1
				mux.Unlock()
			}
		}()
		mux.Lock()
		err2 := fnErr
		mux.Unlock()
		return err2
	})
	wg.Wait()
	if err != nil {
		return err
	}
	mux.Lock()
	err2 := fnErr
	mux.Unlock()
	return err2
}

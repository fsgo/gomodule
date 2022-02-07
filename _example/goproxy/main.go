// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/30

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fsgo/gomodule"
)

func main() {
	gp := &gomodule.GoProxy{
		Module: "github.com/fsgo/fspool",
		Proxy:  "https://goproxy.cn/",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	latest, err := gp.Latest(ctx)
	fmt.Println("Latest:", latest, err)

	versions, err := gp.VersionList(ctx)
	fmt.Printf("VersionList: %q %v\n", versions, err)

	minfo, err := gp.VersionInfo(ctx, "master")
	fmt.Println("VersionInfo(master):", minfo, err)

	info999, err := gp.VersionInfo(ctx, "v0.1.999")
	fmt.Println("VersionInfo(not_found:v0.1.999):", info999, err)

	mod, err := gp.VersionMod(ctx, minfo.Version)
	fmt.Printf("VersionMod(master:%s):\n%s %v", minfo.Version, string(mod), err)

	// zr,err:=gp.VersionZip(ctx,minfo.Version)
	// if err!=nil{
	// 	fmt.Println("VersionZip.err=",err)
	// }else{
	// 	for _,f:=range zr.File{
	// 		fmt.Println(f.Name)
	// 	}
	// }

	zfs, err := gp.VersionFiles(ctx, minfo.Version)
	if err != nil {
		fmt.Println("VersionFiles.err=", err)
	} else {
		for _, zf := range zfs {
			fmt.Println(zf.Name(), zf.IsDir())
		}
	}
}

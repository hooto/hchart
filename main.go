// Copyright 2017 The hchart Authors, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/hooto/hchart/conf"
	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
)

var (
	flagPrefix = flag.String("prefix", "", "the prefix folder path")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	flag.Parse()
	if err := conf.Initialize(*flagPrefix); err != nil {
		hlog.Printf("error", "conf.Initialize error: %v", err)
		os.Exit(1)
	}

	httpsrv.GlobalService.Config.UrlBasePath = "/hchart"
	httpsrv.GlobalService.Config.HttpPort = conf.Config.HttpPort

	module := httpsrv.NewModule("default")

	module.RouteSet(httpsrv.Route{
		Type:       httpsrv.RouteTypeStatic,
		Path:       "~",
		StaticPath: conf.Config.Prefix + "/webui",
	})

	httpsrv.GlobalService.ModuleRegister("/", module)

	fmt.Println("Running")
	httpsrv.GlobalService.Start()
}

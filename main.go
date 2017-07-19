package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/logger"

	"code.hooto.com/hooto/chart/conf"
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
		logger.Printf("error", "conf.Initialize error: %v", err)
		os.Exit(1)
	}

	httpsrv.GlobalService.Config.UrlBasePath = "/chart"
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

package main

import (
	"flag"
	"fmt"

	"github.com/HotPotatoC/kvstore/internal/server"
	"github.com/HotPotatoC/kvstore/internal/version"
	"github.com/HotPotatoC/kvstore/pkg/logger"

	"net/http"
	_ "net/http/pprof"
)

var (
	host  = flag.String("host", "0.0.0.0", "KVStore server host")
	port  = flag.Int("port", 7275, "KVStore server port")
	debug = flag.Bool("debug", false, "Debug mode")
)

func init() {
	flag.StringVar(host, "h", "0.0.0.0", "KVStore server host")
	flag.IntVar(port, "p", 7275, "KVStore server port")
	flag.BoolVar(debug, "d", false, "Debug mode")
}

func main() {
	flag.Parse()
	log := logger.New()
	server := server.New(version.Version, version.Build)

	if *debug {
		log.Info("-=-=-=-=-=-= Running in debug mode =-=-=-=-=-=-")
		go func() {
			log.Infof("Pprof started -> http://%s:%d/debug/pprof", *host, *port+1)
			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port+1), nil); err != nil {
				log.Fatalf("pprof failed: %v", err)
			}
		}()
	}

	log.Fatal(server.Start(*host, *port))
}

package main

import (
	"./server"
	"flag"
	"github.com/eyedeekay/sam-forwarder/config"
)

var (
	manipulateKey = flag.String("b64", "b64", "base64 address to add or remove")
	peerHost      = flag.String("peer", "/etc/i2pd/tunnels.conf.d/i2pd-copy-id.conf", "use this config file")
	samHost       = flag.String("host", "127.0.0.1", "SAM host to use")
	samPort       = flag.String("port", "7656", "SAM port to use")
	removeKey     = flag.Bool("rem", false, "Remove the b64 key")
	server        = flag.Bool("server", false, "Run a service")
)

func main() {
	conf, err := i2ptunconf.NewI2PTunConf(*peerHost)
	if err != nil {
		panic(err)
	}
	if *server {
		serve, err := idserver.NewIDServer(idserver.Config(conf), idserver.ConfigPath(*peerHost))
		if err != nil {
			panic(err)
		}
		if _, err := serve.ListenAndServe(); err != nil {
			panic(err)
		}
	}
}

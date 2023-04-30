package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Vacym/neighbors-force/internal/apiserver"
	"github.com/Vacym/neighbors-force/internal/htmlserver"
	"github.com/Vacym/neighbors-force/internal/proxyserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to config file")
}

func main() {

	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	go startHTMLServer(config)
	fmt.Println("start html server")

	go startAPIServer(config)
	fmt.Println("start api server")

	startProxyServer(config)
	fmt.Println("start proxy server")
}

func startHTMLServer(config *apiserver.Config) {
	if err := htmlserver.Start(":8082"); err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
}

func startAPIServer(config *apiserver.Config) {
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
}

func startProxyServer(config *apiserver.Config) {
	if err := proxyserver.Start(); err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
}
